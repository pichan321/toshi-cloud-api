package middlewares

import (
	"errors"
	"file-api/cloud"
	"file-api/handlers"
	"file-api/structs"
	"fmt"

	"github.com/labstack/echo/v4"
)

func CheckFileHandle(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		handle := c.Param("fileUuid")

		if handle == "" {
			return handlers.ErrorHandler(c, 400, errors.New("bad request"))
		}

		db, err := cloud.GetPostgres()
		if err != nil {
			return handlers.ErrorHandler(c, 500, errors.New("internal Server Error"))
		}

		row := db.QueryRowx(fmt.Sprintf(`SELECT * FROM files WHERE uuid = '%s'`, handle))
		var file structs.File
		err = row.StructScan(&file)

		if err != nil {
			return handlers.ErrorHandler(c, 404, errors.New("requested handle does not exist"))
		}

		defer db.Close()
		return next(c)
	}
}