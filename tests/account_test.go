package tests

import (
	"encoding/json"
	"file-api/cloud"
	"file-api/handlers"
	"file-api/structs"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	user := structs.Account{
		Username: "test",
		Password: "test",
	}
	assert.Equal(t, "", user.Token)
	userJson, err := json.Marshal(user)
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")
	err = handlers.Login(c)
	assert.NoError(t, err)

	err = json.Unmarshal([]byte(rec.Body.String()), &user)
	assert.NoError(t, err)
	assert.Equal(t, "68906ab9-a91b-40f0-ad2a-5f76b31aa734", user.Uuid)
	assert.Equal(t, "test", user.Username)
	assert.NotEqual(t, "", user.Token)
	assert.Equal(t, 200, c.Response().Status)
}

func TestRegisterWithExistingAccountDetails(t *testing.T) {
	user := structs.Account{
		Email: "admin@admin.com",
		Username: "test",
		Password: "admin",
	}
	assert.Equal(t, "", user.Token)
	userJson, err := json.Marshal(user)
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/register")
	err = handlers.RegisterAccount(c)
	assert.NoError(t, err)

	assert.Equal(t, "{\"message\":\"Bad Request\",\"code\":400,\"error\":\"Cannot register account!\"}\n", rec.Body.String())
	assert.Equal(t, 400, c.Response().Status)
}


func TestRegisterNewAccount(t *testing.T) {
	user := structs.Account{
		Email: "admin@admin.com",
		Username: "admin",
		Password: "admin",
	}
	assert.Equal(t, "", user.Token)
	userJson, err := json.Marshal(user)
	assert.NoError(t, err)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/register")
	err = handlers.RegisterAccount(c)
	assert.NoError(t, err)

	db, err := cloud.GetPostgres()
	assert.NoError(t, err)
	defer db.Exec(fmt.Sprintf(`DELETE FROM accounts WHERE username = '%s'`, user.Username))
	assert.Equal(t, "Account successfully registered!", rec.Body.String())
	assert.Equal(t, 200, c.Response().Status)
}