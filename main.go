package main

import (
	"file-api/handlers"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("%v", err)
	}

	router := echo.New()

	//Middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.Recover())
	//router.Use(middleware.Logger())
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	  }))
	
	//GET
	router.GET("/", func (c echo.Context) error {
		return c.String(200, "TOSHI CLOUD")
	})
	router.GET("/get-files/:user", handlers.GetFiles)
	router.GET("/download/:fileUuid", handlers.DownloadFile)
	router.GET("/hide/:fileUuid", handlers.HideFile)
	router.GET("/unhide/:fileUuid", handlers.UnhideFile)
	router.GET("/stream/:fileUuid", handlers.StreamFile)
	router.GET("/content/:fileUuid", handlers.GetFileContent)
	router.GET("/delete/:fileUuid", handlers.DeleteFile)
	
	//POST
	router.POST("/upload", handlers.UploadFile)
	router.POST("/prepare-multipart-upload", handlers.PrepareMultipartUpload)
	router.POST("/multipart-upload", handlers.MultipartUploadFile)
	router.POST("/login", handlers.Login)
	router.POST("/register-account", handlers.RegisterAccount)
	router.POST("/parse-and-upload", handlers.ParseAndUpload)
	
	router.Logger.Fatal(router.Start(":8080"))
}