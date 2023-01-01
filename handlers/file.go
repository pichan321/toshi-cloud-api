package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"file-api/cloud"
	"file-api/structs"
	"file-api/utils"

	"github.com/labstack/echo/v4"
)


func UploadFile(c echo.Context) (err error) {
	db, err := cloud.GetPostgres()
	if err != nil {
		log.Printf("%v", err)
		return ErrorHandler(c, 500, err)
	}
	
	user := c.FormValue("userUuid")
	file, err := c.FormFile("file")
	name := c.FormValue("name")
	size := c.FormValue("size")
	sizeMb := c.FormValue("sizeMb")
	actualSize, _ := strconv.ParseFloat(sizeMb, 32)
	timestamp := time.Now().Format("2006-01-02 15:04:05 PM")
	bucket := utils.GetBucketUuid(actualSize)

	if bucket.Uuid == "" || bucket.AccessToken == "" {
		return ErrorHandler(c, 500, err)
	}
	
	 fileInfo := structs.File{
	 	Uuid: utils.GenerateUuid(),
	 	Name: utils.FixEscape(name),
	 	Size: size,
	 	SizeMb: actualSize,
	 	UploadedDate: timestamp,
	 	UserUuid: user,
	 	BucketUuid: bucket.Uuid,
	 }

	src, err := file.Open()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	ctx := context.Background()
	project, err := cloud.GetStorj(bucket.AccessToken, ctx)
	if err != nil {
		fmt.Println("Errrrrrrrrrrr")
		log.Printf("could not open project: %v", err)
		return ErrorHandler(c, 500, err)
	}
	_, err = project.EnsureBucket(context.Background(), bucket.Name)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}
	storjFilename := utils.StorjFilename(fileInfo.Uuid, fileInfo.Name, "___")
	upload, err := project.UploadObject(ctx, bucket.Name, storjFilename, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}
	data, err := ioutil.ReadAll(src)
	_, err = db.Exec(fmt.Sprintf(`insert into files (uuid, name, size, size_mb, uploaded_date, account_uuid, bucket_uuid, status) values ('%s', '%s', '%s', '%f','%s', '%s', '%s', '1')`, fileInfo.Uuid, storjFilename, fileInfo.Size, fileInfo.SizeMb, fileInfo.UploadedDate, fileInfo.UserUuid, fileInfo.BucketUuid))
	
	if err != nil {
		return fmt.Errorf("could not upload data: %v", err)
	}
	// Copy the data to the upload.
	buf := bytes.NewBuffer(data)
	_, err = io.Copy(upload, buf)
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %v", err)
	}

	// Commit the uploaded object.
	err = upload.Commit()
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf(`update files set status = '100.0' where uuid = '%s'`, fileInfo.Uuid))
	if err != nil {
	 	log.Printf("%v", err)
	 	return ErrorHandler(c, 500, err)
	}
	err = utils.UpdateBucketSize(fileInfo.BucketUuid, fileInfo.SizeMb)
	fmt.Println("Check")
	fmt.Println(fileInfo.BucketUuid)
	fmt.Println(fileInfo.SizeMb)
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	defer db.Close()
	defer project.Close()
	defer src.Close()
	return c.JSON(http.StatusOK, structs.Message{Message: "Uploaded successfully!", Code: 200})
}

func PrepareMultipartUpload(c echo.Context) (err error) {
	db, err := cloud.GetPostgres()
	if err != nil {
		log.Printf("%v", err)
		return ErrorHandler(c, 500, err)
	}

	user := c.FormValue("userUuid")
	name := c.FormValue("name")
	size := c.FormValue("size")
	sizeMb := c.FormValue("sizeMb")
	actualSize, _ :=  strconv.ParseFloat(sizeMb, 32)
	actualSize = math.Floor(actualSize*100)/100
	bucket := utils.GetBucketUuid(actualSize)
	timestamp := time.Now().Format("2006-01-02 15:04:05 PM")


	fileInfo := structs.File{
		Uuid: utils.GenerateUuid(),
		Name: utils.FixEscape(name),
		Size: size,
		SizeMb: actualSize,
		UploadedDate: timestamp,
		UserUuid: user,
		BucketUuid: bucket.Uuid,
	}


	ctx := context.Background()
	project, err := cloud.GetStorj(bucket.AccessToken, ctx)
	if err != nil {
		fmt.Println("Errrrrrrrrrrr multi")
		
		return ErrorHandler(c, 500, err)
	}

	_, err = project.EnsureBucket(ctx, bucket.Name)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}
	storjFilename := utils.StorjFilename(fileInfo.Uuid, fileInfo.Name, "___")
	begin, _ := project.BeginUpload(ctx, bucket.Name, utils.FixEscape(storjFilename), nil)

	db.Exec(fmt.Sprintf(`insert into files (uuid, name, size, size_mb, uploaded_date, account_uuid, bucket_uuid, status, uploadId) values ('%s', '%s', '%s', '%f','%s', '%s', '%s', '%s', '%s')`, fileInfo.Uuid, utils.FixEscape(storjFilename), fileInfo.Size, fileInfo.SizeMb, fileInfo.UploadedDate, fileInfo.UserUuid, fileInfo.BucketUuid, "1.0", begin.UploadID))

	defer db.Close()
	defer project.Close()
	
	return c.JSON(http.StatusOK, structs.Message{Message: begin.UploadID, Code:200, Name: fileInfo.Name})
}

func MultipartUploadFile(c echo.Context) (err error) {
	db, err := cloud.GetPostgres()
	if err != nil {
		log.Printf("%v", err)
		return ErrorHandler(c, 500, err)
	}
	
	// user := c.FormValue("userUuid")
	file, err := c.FormFile("file")
	
	// size := c.FormValue("size")
	//part := c.FormValue("part")
	sizeMb := c.FormValue("sizeMb")
	uploadId := c.FormValue("uploadId")
	actualSize, _ :=  strconv.ParseFloat(sizeMb, 32)
	actualSize = math.Floor(actualSize*100)/100
	current := c.FormValue("current")
	currentPart, _ := strconv.ParseInt(current, 10, 64)
	total := c.FormValue("total")
	totalPart, _ := strconv.ParseInt(total, 10, 64)
	fmt.Println(actualSize)
	// timestamp := time.Now().Format("2006-01-02 15:04:05 PM")
	bucket := utils.GetBucketUuid(actualSize)
	fmt.Printf("%s", "BUCKET INFO")
	fmt.Printf("%v", bucket.Name)
	fmt.Printf("%v", bucket.AccessToken)

	// if bucket.Uuid == "" || bucket.AccessToken == "" {
	// 	return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	// }
	var filename string
	row  := db.QueryRowx(fmt.Sprintf("select name from files where uploadid = '%s'", uploadId))
	row.Scan(&filename)
	fmt.Println(filename)
	// fileInfo := structs.File{
	// 	Uuid: utils.GenerateUuid(),
	// 	Name: name,
	// 	Size: size,
	// 	SizeMb: actualSize,
	// 	UploadedDate: timestamp,
	// 	UserUuid: user,
	// 	BucketUuid: bucket.Uuid,
	// 	Part: currentPart,
	// 	Total: totalPart,
	// }

	src, err := file.Open()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	ctx := context.Background()
	project, err := cloud.GetStorj(bucket.AccessToken, ctx)
	//project, err := cloud.GetStorj(bucket.AccessToken, ctx)

	if err != nil {
		fmt.Println("Couldnt open project")
		log.Printf("could not open project: %v", err)
		return ErrorHandler(c, 500, err)
	}

	_, err = project.EnsureBucket(ctx, bucket.Name)// bucket.Name
	if err != nil {
		fmt.Printf("could not initiate upload: %v", err)
		return ErrorHandler(c, 500, err)
	}

	upload, err := project.UploadPart(ctx, bucket.Name, filename, uploadId, uint32(currentPart)) //uint32(fileInfo.Part)
	if err != nil {
		fmt.Printf("could not initiate upload: %v", err)
		return ErrorHandler(c, 500, err)
	}

	data, err := ioutil.ReadAll(src)
	if err != nil {
		fmt.Printf("could not initiate upload: %v", err)
		return ErrorHandler(c, 500, err)
	}
	// Copy the data to the upload.

	buf := bytes.NewBuffer(data)

	_, err = io.Copy(upload, buf)

	if err != nil {
		_ = upload.Abort()
		fmt.Printf("could not upload data: %v", err)
	}

	// Commit the uploaded object.
	err = upload.Commit()
	fmt.Println("Upload commited")
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %v", err)
	}

	statusFloat := fmt.Sprintf("%.001f", float64(currentPart)/float64(totalPart)*100.0)
	fmt.Println("Current Part")
	fmt.Println(statusFloat)
	// if (currentPart == 1) {
	// 	_, _ = db.Exec(fmt.Sprintf(`insert into files (uuid, name, size, size_mb, uploaded_date, account_uuid, bucket_uuid, status, uploadId) values ('%s', '%s', '%s', '%f','%s', '%s', '%s', '%s', '%s')`, fileInfo.Uuid, fileInfo.Name, fileInfo.Size, fileInfo.SizeMb, fileInfo.UploadedDate, fileInfo.UserUuid, fileInfo.BucketUuid, "1.0", uploadId))
	// 	fmt.Println(fmt.Sprintf(`insert into files (uuid, name, size, size_mb, uploaded_date, account_uuid, bucket_uuid, status, uploadId) values ('%s', '%s', '%s', '%f','%s', '%s', '%s', '%s', '%s')`, fileInfo.Uuid, fileInfo.Name, fileInfo.Size, fileInfo.SizeMb, fileInfo.UploadedDate, fileInfo.UserUuid, fileInfo.BucketUuid, statusFloat, uploadId))
	// } else {
	_, _ = db.Exec(fmt.Sprintf(`update files set status = '%s' where uploadId = '%s'`, statusFloat, uploadId))
	fmt.Println(fmt.Sprintf(`update files set status = '%s' where uploadId = '%s'`, statusFloat, uploadId))

	//  _, err = db.Exec(fmt.Sprintf(`insert into files (uuid, name, size, size_mb, uploaded_date, account_uuid, bucket_uuid) values ('%s', '%s', '%s', '%f','%s', '%s', '%s')`, fileInfo.Uuid, fileInfo.Name, fileInfo.Size, fileInfo.SizeMb, fileInfo.UploadedDate, fileInfo.UserUuid, fileInfo.BucketUuid))

	// if err != nil {
	//  	log.Printf("%v", err)
	//  	return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	// }


	defer db.Close()
	fmt.Println("SAFE UPLOADED PART")
	if (current == total) {
		fmt.Println("CURRENT")
		fmt.Println(current)
		fmt.Println("TOTAL")
		fmt.Println(total)
		defer project.CommitUpload(ctx, bucket.Name, filename, uploadId, nil)
		fmt.Println(filename)
		err = utils.UpdateBucketSize(bucket.Uuid, actualSize)
		if err != nil {
			return ErrorHandler(c, 500, err)
		}
	}
	

	defer project.Close()
	defer src.Close()
	return c.JSON(http.StatusOK, structs.Message{Message: "Uploaded successfully!", Code:200})

}

func DownloadFile(c echo.Context) (err error) {
	db, _ := cloud.GetPostgres()
	fileUuid := c.Param("fileUuid")

	query := fmt.Sprintf(`select access_token, buckets.name as bucket_name, files.name as file_name, buckets.sharelink as share_link from (select * from files where files.uuid = '%s') as files join buckets on files.bucket_uuid = buckets.uuid`, fileUuid) 

	row := db.QueryRowx(query)
	columnNames, _ := row.Columns()

	data := utils.ScanToMap(columnNames, row)

	if fileUuid == "" {
		return ErrorHandler(c, 404, err)
	}

	if data["share_link"] == "" || data["file_name"] == "" {
		return ErrorHandlerWithMsg(c, 404, nil, "Could not get download link")
	}
	fmt.Println("Download File: ")
	fmt.Println(data["share_link"] + "/" + utils.FixEscape(data["file_name"]) + "?wrap=0")

	return c.JSON(http.StatusOK, structs.Message{Message: data["share_link"] + "/" + utils.FixEscape(data["file_name"]) + "?wrap=0", Code: 200})
	
}

func DownloadFileStream(c echo.Context) (err error) {
	db, _ := cloud.GetPostgres()
	bucketUuid := c.Param("bucketUuid")
	fileUuid := c.Param("fileUuid")

	query := fmt.Sprintf(`select access_token, buckets.name as bucket_name, files.name as file_name from (select * from files where files.uuid = '%s') as files join buckets on files.bucket_uuid = buckets.uuid`, fileUuid) 

	row := db.QueryRowx(query)
	columnNames, _ := row.Columns()

	data := utils.ScanToMap(columnNames, row)

	if bucketUuid == "" {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404"})
	}
	if fileUuid == "" {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404"})
	}

	ctx := context.Background()
	project, err := cloud.GetStorj(data["access_token"], ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500", Code: 500})
	}

	download, err := project.DownloadObject(ctx, data["bucket_name"], data["file_name"], nil)
	if err != nil {
		return fmt.Errorf("could not open object: %v", err)
	}

	fmt.Println("Download Link")
	fmt.Printf("%v", download.Info())
	defer download.Close()
	//	receivedContents, err := io.ReadAll(download)
	// Read everything from the download stream
	filename := data["file_name"]
	buf := make([]byte, 64 * 1024)

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, filename))
	c.Response().WriteHeader(http.StatusOK)
	c.Attachment(filename, filename)
    for {
        n, err := io.ReadFull(download, buf)
	if err == io.EOF {
		break
	}
	if err != nil {
		fmt.Println(err)
		continue
	}
	// downlaodedFile.Write(buf[:n])
		a := bytes.NewBuffer(buf[:n])
		io.Copy(c.Response(), a)
    }



	// //fileBuffer := bytes.NewBuffer(receivedContents)
	// if err != nil {
	// 	return fmt.Errorf("could not read data: %v", err)
	// }

	// filename := data["file_name"]
	// downlaodedFile, err := os.Create(filename)
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// }
	// downlaodedFile.Write(fileBuffer.Bytes())

	defer project.Close()
	defer os.Remove(filename)
	// return c.Attachment(filename, filename)
	return nil
}

func GetFiles(c echo.Context) (err error) {
	db, err := cloud.GetPostgres()
	if err != nil {
		log.Printf("%v", err)
		return ErrorHandler(c, 500, err)
	}

	user := c.Param("user")

	if user == "" {
		return ErrorHandler(c, 404, nil)
	}

	query := fmt.Sprintf("select * from files where account_uuid = '%s'", user)

	rows, _ := db.Queryx(query)
	files := []structs.File{}
	
	for rows.Next() {
		file := structs.File{}
		err = rows.StructScan(&file)
		if err != nil {
			fmt.Println(err)
		}
		files = append(files, file)

	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].UploadedDate > files[j].UploadedDate
	})
	

	defer db.Close()
	return c.JSON(http.StatusOK, files)
}

func DeleteFile(c echo.Context) (err error) {
	db, err := cloud.GetPostgres()
	fileUuid := c.Param("fileUuid")

	query := fmt.Sprintf(`select access_token, buckets.name as bucket_name, buckets.uuid as bucket_uuid, files.name as file_name, files.size_mb as file_size from (select * from files where files.uuid = '%s') as files join buckets on files.bucket_uuid = buckets.uuid`, fileUuid) 
	row := db.QueryRowx(query)
	columnNames, _ := row.Columns()
	data := utils.ScanToMap(columnNames, row)

	ctx := context.Background()
	project, err := cloud.GetStorj(data["access_token"], ctx)
	if err != nil {
		return fmt.Errorf("could not request access grant: %v", err)
	}

	_, err = project.DeleteObject(ctx, data["bucket_name"], utils.FixEscape(data["file_name"]))
	
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	_, err = db.Exec(fmt.Sprintf(`delete from files where uuid = '%s'`, fileUuid))

	if err != nil {
		return ErrorHandler(c, 500, err)
	}
	fileSize, _ := strconv.ParseFloat(data["file_size"], 32)
	err = utils.UpdateBucketSize(data["bucket_uuid"], -fileSize)
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	defer db.Close()
	defer project.Close()
	return c.JSON(http.StatusOK, structs.Message{Message: "File deleted successfully!", Code: 200})
}

func StreamFile(c echo.Context) (err error) {
	db, err := cloud.GetPostgres()

	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	// bucketUuid := c.Param("bucketUuid")
	fileUuid := c.Param("fileUuid")
	fmt.Println(fileUuid)
	query := fmt.Sprintf(`select access_token, buckets.name as bucket_name, files.name as file_name, buckets.shareLink as share_link  from (select * from files where files.uuid = '%s') as files join buckets on files.bucket_uuid = buckets.uuid`, fileUuid) 

	row := db.QueryRowx(query)
	columnNames, _ := row.Columns()

	data := utils.ScanToMap(columnNames, row)
	fmt.Println(data)
	fmt.Println(fmt.Sprintf(`%s/%s?wrap=0`, data["share_link"], data["file_name"]))
	return c.JSON(http.StatusOK, structs.Message{Message: fmt.Sprintf(`%s/%s?wrap=0`, data["share_link"], data["file_name"]), Code: 200})
}

func GetFileContent(c echo.Context) (err error) {
	db, _ := cloud.GetPostgres()
	fileUuid := c.Param("fileUuid")

	query := fmt.Sprintf(`select access_token, buckets.name as bucket_name, files.name as file_name from (select * from files where files.uuid = '%s') as files join buckets on files.bucket_uuid = buckets.uuid`, fileUuid) 

	row := db.QueryRowx(query)
	columnNames, _ := row.Columns()

	data := utils.ScanToMap(columnNames, row)

	if fileUuid == "" {
		return ErrorHandler(c, 404, err)
	}

	ctx := context.Background()
	project, err := cloud.GetStorj(data["access_token"], ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500", Code: 500})
	}

	download, err := project.DownloadObject(ctx, data["bucket_name"], data["file_name"], nil)
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	defer download.Close()

	receivedContents, err := io.ReadAll(download)

	fileType := strings.Split(data["file_name"], ".")[1]

	if fileType == "txt" {
		return c.Blob(http.StatusOK, "text/plain", receivedContents)
	}

	imagesSlice := []string{"jpg", "png", "gif"} 
	var imagesInterfaceSlice []interface{}
	for _, v := range imagesSlice {
		imagesInterfaceSlice = append(imagesInterfaceSlice, v)
	}
	

	imageType := utils.ExistsWithin(imagesInterfaceSlice, fileType)
	if  imageType != "" {
		return c.Blob(http.StatusOK, fmt.Sprintf(`image/%s`, imageType), receivedContents)
	}

	return nil

}

func HideFile(c echo.Context) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}
	
	fileUuid := c.Param("fileUuid")
	if fileUuid == "" {
		return ErrorHandlerWithMsg(c, 404, err, "The file does not exist.")
	}

	_, err = db.Exec(fmt.Sprintf(`UPDATE files SET hidden = true where uuid = '%s'`, fileUuid))
	if err != nil {
		return ErrorHandlerWithMsg(c, 404, err, "Could not hide the provided file.")
	}
	defer db.Close()

	return Success(c, "File is hidden")
}


func UnhideFile(c echo.Context) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}
	
	fileUuid := c.Param("fileUuid")
	if fileUuid == "" {
		return ErrorHandlerWithMsg(c, 404, err, "The file does not exist.")
	}

	_, err = db.Exec(fmt.Sprintf(`UPDATE files SET hidden = false where uuid = '%s'`, fileUuid))
	if err != nil {
		return ErrorHandlerWithMsg(c, 404, err, "Could not hide the provided file.")
	}
	defer db.Close()

	return Success(c, "File is unhidden")
}