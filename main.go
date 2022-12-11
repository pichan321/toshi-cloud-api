package main

import (
	"file-api/handlers"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//var env map[string]string, _ := godotenv.Read(".env")
const BUCKET_SIZE_LIMIT = 150000.00

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("%v", err)
	}

	router := echo.New()

	router.Use(middleware.CORS())
	router.Use(middleware.Recover())
	//router.Use(middleware.Logger())
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	  }))
	
	//GET

	router.GET("/", func (c echo.Context) error {
		return c.String(200, os.Getenv("ACCESS_KEY"))
	})
	router.GET("/get-files/:user", handlers.GetFiles)
	router.GET("/download/:user/:handle", handlers.DownloadFile)
	
	//POST
	router.POST("/upload", handlers.UploadFile)
	router.POST("/login", handlers.Login)
	router.POST("/register-account", handlers.RegisterAccount)
	router.Logger.Fatal(router.Start(":8080"))
}