package handlers

import (
	"errors"
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
		return ErrorHandler(c, 500, err)
	}

	var account structs.Account
	err = c.Bind(&account)
	if err != nil {
		return ErrorHandler(c, 400, err)
	}

	err = utils.ValidateEmail(account.Email) 
	if err != nil {
		return ErrorHandler(c, 400, errors.New("invalid email address"))
	}

	if (account.Username == "" || account.Email == "" || account.Password == "") {
		return ErrorHandler(c, 400, errors.New("bad request"))
	}

	var checkAccount structs.Account
	accounts := db.QueryRowx(fmt.Sprintf(`select * from accounts where username = '%s' or email = '%s'`, account.Username, account.Email))
	accounts.StructScan(&checkAccount)

	if (account.Username == checkAccount.Username) || (account.Email == checkAccount.Email) {
		return ErrorHandler(c, 400, errors.New("Cannot register account!"))
	}

	id := uuid.New()
	hashedPassword := utils.HashPassword(account.Password)

	_, err = db.Exec(fmt.Sprintf(`insert into accounts (uuid, username, password, email) values ('%s', '%s','%s','%s')`,  id.String(), strings.ToLower(account.Username), string(hashedPassword), account.Email))
	if err != nil {
		fmt.Printf("%v", err)
	}
	defer db.Close()
	return c.JSON(http.StatusOK, structs.Message{Code: 200})
}

func Login(c echo.Context) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return ErrorHandler(c, 500, err)
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
		return ErrorHandler(c, 404, err)
	}

	row := db.QueryRowx(fmt.Sprintf(`select * from accounts where username = '%s' limit 1`, account.Username))

	var dbAccount structs.Account
	row.StructScan(&dbAccount)

	if (utils.HashPassword(account.Password) != dbAccount.Password || (strings.Compare(account.Username, dbAccount.Username) != 0)) {
		return ErrorHandler(c, 404, err)
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
	return c.JSON(http.StatusOK, dbAccount)
}