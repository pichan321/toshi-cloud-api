package handlers

import (
	"file-api/structs"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Success(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, structs.Message{Message: message, Code: 200, Error: "none"})
}

func ErrorHandler(c echo.Context, code int, err error) error {
	switch (code) {
		case 400: return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request", Code: 400, Error: err.Error()})
		case 404: return c.JSON(http.StatusNotFound, structs.Message{Message: "Bad Request", Code: 404, Error: err.Error()})
		case 500: return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error", Code: 500, Error: err.Error()})
	}
	return nil
}

func ErrorHandlerWithMsg(c echo.Context, code int, err error, message string) error {
	switch (code) {
		case 400: return c.JSON(http.StatusBadRequest, structs.Message{Message: message, Code: 400, Error: fmt.Sprintf("%s", err)})
		case 404: return c.JSON(http.StatusNotFound, structs.Message{Message: message, Code: 404, Error: fmt.Sprintf("%s", err)})
		case 500: return c.JSON(http.StatusInternalServerError, structs.Message{Message: message, Code: 500, Error: fmt.Sprintf("%s", err)})
	}
	return nil
}
