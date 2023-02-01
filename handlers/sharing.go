package handlers

import (
	"errors"
	"file-api/cloud"
	"file-api/structs"
	"file-api/utils"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SharedWith(c echo.Context) error {

	return nil
}

func ShareFileShareAccess(c echo.Context) error {
	db, err := cloud.GetPostgres()

	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	var share structs.ShareFile
	c.Bind(&share)

	query := fmt.Sprintf(`SELECT count(*) FROM sharing where handle = '%s' and file_owner = '%s' and file_recipient = '%s'`, share.Handle, share.Owner, share.Recipient)
	var count int
	row := db.QueryRowx(query)
	row.Scan(&count)

	if count > 0 {
		return ErrorHandler(c, 400, errors.New("file is already shared with the recipient"))
	}

	query = fmt.Sprintf(`INSERT INTO sharing VALUES ('%s', '%s', '%s', '%s')`, utils.GenerateUuid(), share.Handle, share.Owner, share.Recipient)

	db.Exec(query)
	
	defer db.Close()
	return c.JSON(http.StatusOK, share)
}

func ShareFileRevokeAccess(c echo.Context) error {
	db, err := cloud.GetPostgres()

	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	var share structs.ShareFile
	c.Bind(&share)

	query := fmt.Sprintf(`DELETE FROM sharing WHERE file_owner = '%s' and file_recipient = '%s' and handle = '%s'`, share.Owner, share.Recipient, share.Handle)

	db.Exec(query)
	
	defer db.Close()
	return c.JSON(http.StatusOK, fmt.Sprintf(`%s's access to your file has been revoked.`, share.Recipient))
}

func DeleteSharedFile(c echo.Context) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}
	var share structs.ShareFile
	c.Bind(&share)
	fmt.Println("DELETE SHARE")
	fmt.Println(share)

	db.Exec(fmt.Sprintf(`DELETE FROM sharing WHERE handle = '%s' and file_owner = '%s' and file_recipient = '%s'`, share.Handle, share.Owner, share.Recipient))

	defer db.Close()
	return c.JSON(http.StatusOK, "DELETED")
}

func GetUsersToShare(c echo.Context) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}
	user := c.Param("user")
	username := c.Param("username")
	handle := c.Param("handle")
	query := fmt.Sprintf("select uuid, username from accounts where  username ILIKE '%%%s%%' and uuid != '%s'", username, user)
	query = fmt.Sprintf(`select t1.uuid, t1.username, t2.handle, '' as shared from (select uuid, username from accounts where  username ILIKE '%%%s%%' and uuid != '%s') as t1 left join (select * from sharing where handle = '%s') as t2 on t1.uuid = t2.file_recipient`, username, user, handle)
	fmt.Println("SHARED")
	fmt.Println(query)
	rows, err := db.Queryx(query)

	colNames, _ := rows.Columns()
	// users := []structs.SharedUser{}
	users := []map[string]interface{}{}
	for rows.Next() {
		user := utils.ScanToMapRows(colNames, rows)
		if user["handle"] == "" {
			user["shared"] = false
			users = append(users, user)
			continue
		}
		user["shared"] = true
		users = append(users, user)
	
	}
	return c.JSON(http.StatusOK, users)
}