package tests

import (
	"file-api/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

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

func TestDownloadFileExists(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/download/:fileUuid")
	c.SetParamNames("fileUuid")
	c.SetParamValues("1724a4ce-3e1b-4562-abb6-bb166a746414")
	err := handlers.DownloadFile(c)
	assert.NoError(t, err)
	assert.Equal(t, "{\"message\":\"https://link.storjshare.io/jvu3hs64ove55hhtca7ngpbrybzq/vattana/1724a4ce-3e1b-4562-abb6-bb166a746414___test.txt?wrap=0\",\"code\":200}\n", rec.Body.String())
	assert.Equal(t, 200, c.Response().Status)
}

func TestDownladFileNotExist(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/download/:fileUuid")
	c.SetParamNames("fileUuid")
	c.SetParamValues("83487944-a624-4dfb-9c80-45a6596244b1")
	err := handlers.DownloadFile(c)
	assert.NoError(t, err)
	assert.Equal(t, "{\"message\":\"Could not get download link\",\"code\":404,\"error\":\"%!s(\\u003cnil\\u003e)\"}\n", rec.Body.String())
	assert.Equal(t, 404, c.Response().Status)
}

func TestHideFile(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/hide/:fileUuid")
	c.SetParamNames("fileUuid")
	c.SetParamValues("1724a4ce-3e1b-4562-abb6-bb166a746414")
	err := handlers.HideFile(c)
	assert.NoError(t, err)
	assert.Equal(t, "{\"message\":\"File is hidden\",\"code\":200,\"error\":\"none\"}\n", rec.Body.String())
	assert.Equal(t, 200, c.Response().Status)
}

func TestUnhideFile(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/unhide/:fileUuid")
	c.SetParamNames("fileUuid")
	c.SetParamValues("1724a4ce-3e1b-4562-abb6-bb166a746414")
	err := handlers.UnhideFile(c)
	assert.NoError(t, err)
	assert.Equal(t, "{\"message\":\"File is unhidden\",\"code\":200,\"error\":\"none\"}\n", rec.Body.String())
	assert.Equal(t, 200, c.Response().Status)
}