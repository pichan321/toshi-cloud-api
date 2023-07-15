package middlewares

import (
	"context"
	"errors"
	"file-api/cloud"
	"file-api/handlers"
	"file-api/structs"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func doesUserExistAndInit(db *sqlx.DB, userClaims structs.CustomClaims) string {
	row := db.QueryRowx(fmt.Sprintf(`SELECT * FROM accounts WHERE accounts.auth0 = '%s'`, userClaims.Sub))
	var account structs.Account
	err := row.StructScan(&account)
	if err != nil {
		initializeUser(db, userClaims)
		row := db.QueryRowx(fmt.Sprintf(`SELECT * FROM accounts WHERE accounts.auth0 = '%s'`, userClaims.Sub))
		var account structs.Account
		row.StructScan(&account)
		return account.Uuid
	}

	return account.Uuid
}

func initializeUser(db *sqlx.DB, userClaims structs.CustomClaims) {
	_, err := db.Exec(fmt.Sprintf(`INSERT INTO accounts (uuid, username, email, auth0, api_key) VALUES 
	('%s', '%s', '%s', '%s', '')`, uuid.New(), userClaims.Nickname, userClaims.Email, userClaims.Sub))
	if err != nil {
		
	}
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		token := c.Request().Header.Get("authorization")
		values := strings.Split(token, " ")
		if len(values) != 2 {
			return c.JSON(http.StatusBadRequest, "valid token is required")
		}

		tokenType, tokenValue := values[0], values[1]
		if tokenType != "Bearer" {
			return c.JSON(http.StatusBadRequest, "A bearer token is required")
		}

		issuerURL, err := url.Parse("https://dev-qi6tdasmjtbp226c.us.auth0.com/")
		if err != nil {
			log.Fatalf("Failed to parse the issuer url: %v", err)
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		userClaims := structs.CustomClaims{}
		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{"tU3bRZIN5qoka4fS5DRC3CgSlhv93cKw"},
			validator.WithCustomClaims(
				func() validator.CustomClaims {
					return &userClaims
				},
			),
			validator.WithAllowedClockSkew(time.Minute),
		)
		_, err = jwtValidator.ValidateToken(context.TODO(), tokenValue)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error retrieving your information")
		}
		db, err := cloud.GetPostgres()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Internal Server Error")
		}

		userUuid := doesUserExistAndInit(db, userClaims)

		defer db.Close()
		c.Set("user", userClaims)
		c.Set("userUuid", userUuid)
		return next(c)
	}
}

func CheckFileHandle(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		handle := c.Param("fileUuid")

		if handle == "" {
			return handlers.ErrorHandler(c, 404, errors.New("file handle doesn't exist"))
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