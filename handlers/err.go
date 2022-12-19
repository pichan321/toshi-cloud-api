package handlers

import (
	"file-api/structs"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(c echo.Context, code int) error {
	switch (code) {
		case 404: return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Bad Request", Code: 404})
		case 500: return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error", Code: 500})
	}
	return nil
}

func ErrorHandlerWithMsg(c echo.Context, code int, message string) error {
	switch (code) {
		case 404: return c.JSON(http.StatusInternalServerError, structs.Message{Message: message, Code: 404})
		case 500: return c.JSON(http.StatusInternalServerError, structs.Message{Message: message, Code: 500})
	}
	return nil
}
