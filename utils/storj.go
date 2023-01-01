package utils

import (
	"errors"
	"file-api/cloud"
	"fmt"
)

func StorjFilename(uuid string, filename string, delimiter string) string {
	return fmt.Sprintf("%s%s%s", uuid, delimiter, filename)
}

func UpdateBucketSize(uuid string, sizeToUpdate float64) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return errors.New("could not get PostgreSQL")
	}

	var currentBucketSize float64
	row := db.QueryRowx(fmt.Sprintf(`SELECT size FROM buckets where uuid = '%s'`, uuid))
	row.Scan(&currentBucketSize)
	fmt.Println(currentBucketSize)
	newSize := sizeToUpdate + currentBucketSize
	_, err = db.Exec(fmt.Sprintf(`UPDATE buckets SET size = %.2f where uuid = '%s'`, newSize, uuid))
	fmt.Println(fmt.Sprintf(`UPDATE buckets SET size_mb = %.2f where uuid = '%s'`, newSize, uuid))
	if err != nil {
		return errors.New("could not update bucket size")
	}

	defer db.Close()
	return nil
}