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
		return c.JSON(http.StatusInternalServerError, structs.Message{Message: "Internal Server Error 500", Code: 500})
	}

	var account structs.Account
	
	err = c.Bind(&account)
	if err != nil {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404", Code: 404})
	}

	if (account.Username == "" || account.Email == "" || account.Password == "") {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404", Code: 404})
	}

	var checkAccount structs.Account
	accounts := db.QueryRowx(fmt.Sprintf(`select * from accounts where username = '%s'`, account.Username))
	accounts.StructScan(&checkAccount)


	fmt.Printf("%v", checkAccount)
	if (account.Username == checkAccount.Username) {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "Bad Request 404", Code: 404})
	}

	id := uuid.New()
	hashedPassword := utils.HashPassword(account.Password)

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

	token := c.QueryParam("token")
	if (token != "") {
		row := db.QueryRowx(fmt.Sprintf(`select * from accounts where token = '%s'`, token))
		var account structs.Account
		row.StructScan(&account)
		if (account.Token == token) {
			var temp = &account
			temp.Password = ""
			return c.JSON(http.StatusOK, account)
		}
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

	if (utils.HashPassword(account.Password) != dbAccount.Password || (strings.Compare(account.Username, dbAccount.Username) != 0)) {
		return c.JSON(404, structs.Message{Message: "Bad Request 404", Code: 404})
	}

	var tempPointer = &dbAccount
	tempPointer.Password = ""
	tempPointer.Token = utils.GenerateToken()
	//query := fmt.Sprintf(`update accounts set token = '%s' where uuid = '%s'`, dbAccount.Token, dbAccount.Username)
	query := fmt.Sprintf("update accounts set token = '%s' where uuid = '%s'", dbAccount.Token, dbAccount.Uuid)
	_, err = db.Exec(query)
	if err != nil {
		fmt.Printf("%v", err)
	}

	defer db.Close()
	return c.JSON(200, dbAccount)
}