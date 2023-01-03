package main

import (
	"file-api/handlers"
	"file-api/middlewares"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 10, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string,err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}


	err := godotenv.Load()
	if err != nil {
		fmt.Printf("%v", err)
	}

	router := echo.New()

	//Middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.Recover())
	router.Use(middleware.RateLimiterWithConfig(config))
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${status}] uri=${uri} [${method}]\n",
		Output: router.StdLogger.Writer(),
	  }))
	
	//GET
	router.GET("/", func (c echo.Context) error {
		return c.String(200, "TOSHI CLOUD")
	})
	router.GET("/get-files/:user", handlers.GetFiles)

	//POST
	router.POST("/upload", handlers.UploadFile)
	router.POST("/prepare-multipart-upload", handlers.PrepareMultipartUpload)
	router.POST("/multipart-upload", handlers.MultipartUploadFile)
	router.POST("/login", handlers.Login)
	router.POST("/register-account", handlers.RegisterAccount)
	router.POST("/parse-and-upload", handlers.ParseAndUpload)

	//FILE GROUP
	file := router.Group("/file", middlewares.CheckFileHandle)
	file.GET("/download/:fileUuid", handlers.DownloadFile)
	file.GET("/hide/:fileUuid", handlers.HideFile)
	file.GET("/unhide/:fileUuid", handlers.UnhideFile)
	file.GET("/stream/:fileUuid", handlers.StreamFile)
	file.GET("/content/:fileUuid", handlers.GetFileContent)
	file.GET("/delete/:fileUuid", handlers.DeleteFile)
	
	router.Logger.Fatal(router.Start(":8080"))
}