package cloud

import (
	"context"
	"errors"
	"file-api/structs"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"storj.io/uplink"
)

func GetPostgres() (*sqlx.DB, error) {
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s", // sslmode=verify-full
    "pichan-2902.g8z.cockroachlabs.cloud", 26257, "pichan", "OGdtBNEQGFcGS818wuLbxA", "toshi-cloud")
	db, err := sqlx.Open("postgres", postgresqlDbInfo)

	if err != nil {
        return nil, errors.New("could not connect to Postgres")
    }

	return db, nil
}

// func GetMongo() (*mongo.Client, error) {
// 	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
// 	clientOptions := options.Client().
// 	    ApplyURI("mongodb+srv://vattana:1234567890@cluster.sk4xlbi.mongodb.net/?retryWrites=true&w=majority").
// 	    SetServerAPIOptions(serverAPIOptions)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, clientOptions)

// 	if err != nil {
// 	    return nil, errors.New("could not connect to MongoDB")
// 	}
	
// 	return client, nil
// }

func GetStorj(accessToken string, ctx context.Context) (*uplink.Project, error) {
	access, err := uplink.ParseAccess(accessToken)
	if err != nil {
		return nil, errors.New("could not parse project access")
	}

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, errors.New("could not open Storj project")
	}

	return project, nil
}

func GetGorm() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s", // sslmode=verify-full
    "pichan-2902.g8z.cockroachlabs.cloud", 26257, "pichan", "OGdtBNEQGFcGS818wuLbxA", "toshi-cloud")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}
	return db, nil
}

func Migrate(c echo.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s", // sslmode=verify-full
    "pichan-2902.g8z.cockroachlabs.cloud", 26257, "pichan", "OGdtBNEQGFcGS818wuLbxA", "toshi-cloud")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	db.AutoMigrate(&structs.Account{}, &structs.Bucket{}, &structs.File{}, &structs.ShareFile{})
	return c.JSON(http.StatusOK, structs.Message{Code: 200, Message: "Done"})
}
