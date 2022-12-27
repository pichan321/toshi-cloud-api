package tests

import (
	"file-api/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHelloName(t *testing.T) {
	assert.Equal(t, 1, 1, "The two words should be the same.")

}

func TestGetFiles(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/get-files/:user")
	c.SetParamNames("user")
	c.SetParamValues("83487944-a624-4dfb-9c80-45a6596244b1")
	err := handlers.GetFiles(c)
	assert.NoError(t, err)
	assert.Equal(t, 200, c.Response().Status)
}