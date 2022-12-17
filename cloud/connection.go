package cloud

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"storj.io/uplink"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func GetPostgres() (*sqlx.DB, error) {
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=verify-full",
    "pichan-2902.g8z.cockroachlabs.cloud", 26257, "pichan", "OGdtBNEQGFcGS818wuLbxA", "toshi-cloud")
	db, err := sqlx.Open("postgres", postgresqlDbInfo)

	if err != nil {
        return nil, errors.New("could not connect to Postgres")
    }
	
	return db, nil
}

func GetMongo() (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
	    ApplyURI("mongodb+srv://vattana:1234567890@cluster.sk4xlbi.mongodb.net/?retryWrites=true&w=majority").
	    SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
	    return nil, errors.New("could not connect to MongoDB")
	}
	
	return client, nil
}

func GetSQL() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", "1Qstk6LRfj:It0AEZPaCt@tcp(remotemysql.com)/1Qstk6LRfj")

    if err != nil {
        return nil, errors.New("could not connect to MySQL")
    }

	return db, nil
}

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
