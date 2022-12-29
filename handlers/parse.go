package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"file-api/cloud"
	"file-api/structs"
	"file-api/utils"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/csimplestring/go-csv/detector"
	"github.com/labstack/echo/v4"
)

func ParseAndUpload(c echo.Context) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	fileToParse := structs.FileParse{}
	err = c.Bind(&fileToParse)
	if err != nil {
		return ErrorHandler(c, 404, err)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05 PM")

	if fileToParse.Filename == "" {
		temp := &fileToParse
		temp.Filename = timestamp
	}

	bucket := utils.GetBucketUuid(1.0)
	fileInfo := structs.File{
		Uuid: utils.GenerateUuid(),
		Name: utils.FixEscape(fileToParse.Filename),
		Size: "",
		SizeMb: 0,
		UploadedDate: timestamp,
		UserUuid: fileToParse.User,
		BucketUuid: bucket.Uuid,
	}

	ctx := context.Background()
	project, err := cloud.GetStorj(bucket.AccessToken, ctx)
	if err != nil {
		log.Printf("could not open project: %v", err)
		return ErrorHandler(c, 500, err)
	}

	_, err = project.EnsureBucket(context.Background(), bucket.Name)
	if err != nil {
		return ErrorHandler(c, 500, err)
	}
	storjFilename := utils.StorjFilename(fileInfo.Uuid, fileInfo.Name, "___")
	fileType := processParsedText(fileToParse.Content)
	finalFilename := storjFilename + fileType

	upload, err := project.UploadObject(ctx, bucket.Name, finalFilename, nil)
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	_, err = db.Exec(fmt.Sprintf(`insert into files (uuid, name, size, size_mb, uploaded_date, account_uuid, bucket_uuid, status) values ('%s', '%s', '%s', '%f','%s', '%s', '%s', '1')`, fileInfo.Uuid, finalFilename, fileInfo.Size, fileInfo.SizeMb, fileInfo.UploadedDate, fileInfo.UserUuid, fileInfo.BucketUuid))
	
	if err != nil {
		return fmt.Errorf("could not upload data: %v", err)
	}


	buf := bytes.NewBuffer([]byte(fileToParse.Content))
	_, err = io.Copy(upload, buf)

	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %v", err)
	}



	// Commit the uploaded object.
	err = upload.Commit()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}
	_, _ = db.Exec(fmt.Sprintf(`update files set status = '100.0' where uuid = '%s'`, fileInfo.Uuid))

	defer project.Close() 
	defer db.Close()
	return nil
}

func processParsedText(content string) string {
	check := []string{
		isJSON(content),
		isCSV(content),
	}
	fmt.Println("Processed")
	fmt.Printf("%v", check)
	
	for _, v := range check {
		if v != "" {
			return v
		}
	}

	return check[0]
}

func isCSV(content string) string {
	detector := detector.New()
	bytesBuffer := bytes.NewBufferString(content)
	delimiters := detector.DetectDelimiter(bytesBuffer, '"')
	validDelimiters := []rune{',', '\t', '|'}
	var valid bool = false
	if len(delimiters) > 0 {

		for _, v := range validDelimiters {
			if string(v) == delimiters[0] {
				valid = true
			}
		}
		if valid {
			return ".csv"
		} 
	}

	return ""
}

func isJSON(content string) string {
	var js interface{}
	fmt.Println("Detecting JSON")
	if  (json.Unmarshal([]byte(content), &js) == nil) {
		fmt.Println("is JSON file")
		return ".json"
	}
	return ""
}
