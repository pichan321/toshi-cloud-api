package handlers

import (
	"file-api/cloud"
	"file-api/structs"
	"file-api/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterAccount(c echo.Context) error {
	db, err := cloud.GetPostgres()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500"})
	}

	account := new(structs.Account)
	err = c.Bind(&account)
	if err != nil {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404"})
	}

	id := uuid.New()
	hashedPassword := utils.HashPassword(account.Password)
	fmt.Println(hashedPassword)
	_, err = db.Exec(fmt.Sprintf(`insert into accounts (uuid, username, password, email) values ('%s', '%s','%s','%s')`,  id.String(), strings.ToLower(account.Username), string(hashedPassword), account.Email))
	if err != nil {
		fmt.Printf("%v", err)
	}
	defer db.Close()
	return c.JSON(http.StatusOK,  id.String() + "\t" + account.Username + "\t" + hashedPassword + "\t" + account.Email)
}

func Login(c echo.Context) error {
	db, err := cloud.GetPostgres()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500", Code: 500})
	}

	var account structs.Account
	err = c.Bind(&account)

	if err != nil {
		c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404", Code: 404})
	}

	row := db.QueryRowx(fmt.Sprintf(`select * from accounts where username = '%s' limit 1`, account.Username))

	//var dbAccount structs.Account
	var dbAccount structs.Account
	row.StructScan(&dbAccount)


	fmt.Printf("%v", dbAccount)
	fmt.Println()
	fmt.Printf("%v", account)
	fmt.Println(utils.HashPassword(account.Password) == dbAccount.Password)
	fmt.Println(account.Username != dbAccount.Username)
	fmt.Println(account.Username)
	fmt.Println(dbAccount.Username)
	fmt.Println(strings.Compare(account.Username, dbAccount.Username))
	if (utils.HashPassword(account.Password) != dbAccount.Password || (strings.Compare(account.Username, dbAccount.Username) != 0)) {
		return c.JSON(404, structs.Message{Message: "Bad Request 404", Code: 404})
	}

	var a = &dbAccount
	a.Password = ""

	_, err = db.Exec(fmt.Sprintf(`update accounts set token = \'%s\' where uuid = \'%s\'`, utils.GenerateToken(), dbAccount.Uuid))
	if err != nil {
		fmt.Errorf("%v", err)
	}

	defer db.Close()
	return c.JSON(200, dbAccount)
}