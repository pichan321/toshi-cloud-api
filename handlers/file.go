package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"file-api/structs"
	"file-api/utils"

	"github.com/labstack/echo/v4"
	"storj.io/uplink"
)


func UploadFile(c echo.Context) (err error) {
	//db, err := cloud.GetPostgres()
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	}
	
	//user := c.Param("userUuid")
	file, err := c.FormFile("file")
	name := c.FormValue("name")
	//size := c.FormValue("size")
	sizeMb := c.FormValue("sizeMb")
	actualSize, _ := strconv.ParseFloat(sizeMb, 32)
	//timestamp := time.Now().Format("2006-01-02 15:04:05 PM")
	bucket := utils.GetBucketUuid(actualSize)

	if bucket.Uuid == "" || bucket.AccessToken == "" {
		return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	}
	
	// fileInfo := structs.File{
	// 	Uuid: utils.GenerateUuid(),
	// 	Name: name,
	// 	Size: size,
	// 	SizeMb: actualSize,
	// 	UploadedDate: timestamp,
	// 	UserUuid: user,
	// 	BucketUuid: bucket.Uuid,
	// }
	
	// if user == "" {
	// 	log.Printf("%v", err)
	// 	return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404"})
	// }

	// if err != nil {
	// 	log.Printf("%v", err)
	// 	return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404"})
	// }
	
	src, err := file.Open()

	if err != nil {
		return  c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	}

	defer src.Close()
	access, err := uplink.ParseAccess(bucket.AccessToken)
	if err != nil {
		return fmt.Errorf("could not request access grant: %v", err)
	}
	ctx := context.Background()
	// Open up the Project we will be working with.
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(context.Background(), bucket.Name)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}

	// Intitiate the upload of our Object to the specified bucket and key.
	upload, err := project.UploadObject(ctx, bucket.Name, name, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}
	data, err := ioutil.ReadAll(src)

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

	// if err != nil {
	// 	log.Printf("%v", err)
	// 	return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	// }


	// _, err = sql.Exec(fmt.Sprintf(`insert into files (username, handle, filename, filetype, filesize, datestamp) values ("%s", "%s", "%s", "%s","%s", "%s")`, user, filename, filename, filetype, filesize, timestamp))

	// if err != nil {
	// 	log.Printf("%v", err)
	// 	return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	// }

	// if err != nil {
	// 	log.Printf("%v", err)
	// 	return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	// }

	// defer sql.Close()

	return c.JSON(http.StatusOK, structs.Message{Message: "Uploaded successfully!"})

}

func DownloadFile(c echo.Context) (err error) {
	username := c.Param("user")
	handle := c.Param("handle")

	if username == "" {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404"})
	}
	if handle == "" {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404"})
	}

	access, err := uplink.ParseAccess(os.Getenv("ACCESS_KEY"))
	if err != nil {
		return fmt.Errorf("could not request access grant: %v", err)
	}
	ctx := context.Background()

	project, err := uplink.OpenProject(ctx, access)
	download, err := project.DownloadObject(ctx, os.Getenv("BUCKET"), handle, nil)
	if err != nil {
		return fmt.Errorf("could not open object: %v", err)
	}
	defer download.Close()
	// Read everything from the download stream
	receivedContents, err := io.ReadAll(download)
	fileBuffer := bytes.NewBuffer(receivedContents)
	if err != nil {
		return fmt.Errorf("could not read data: %v", err)
	}
	downlaodedFile, err := os.Create("test")
	downlaodedFile.Write(fileBuffer.Bytes())

	defer os.Remove("test")
	return c.Attachment("test", "test")
}	

func GetFiles(c echo.Context) (err error) {
	// sql, err := cloud.GetSQL()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// user := c.Param("user")
	// query := fmt.Sprintf("select * from files where username = \"%s\"", user)
	// rows, err := sql.Query(query)
	// if err != nil {
	// 	panic(err)
	// }

	// files := []structs.File{}
	// for rows.Next() {
	// 	file := structs.File{}
	// 	err = rows.Scan(&file.Id, &file.Username, &file.Handle, &file.FileName, &file.FileType, &file.FileSize, &file.Date)
	// 	fmt.Printf("%v", user)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	files = append(files, file)

	// }

	// defer sql.Close()
	return c.JSON(http.StatusOK, structs.File{})
}

func DeleteFile(c echo.Context) (err error) {

	return nil
}