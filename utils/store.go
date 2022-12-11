package utils

import (
	"file-api/cloud"
	"file-api/structs"
)

const BUCKET_SIZE_LIMIT = 150000.00

func GetBucketUuid(fileSizeMb float64) structs.Bucket {
	db, _ := cloud.GetPostgres()
	defer db.Close()
	
	rows, _ := db.Queryx("select * from buckets")

	for rows.Next(){
		bucket := structs.Bucket{}
		rows.StructScan(&bucket)
		if (bucket.Size + float32(fileSizeMb) <= BUCKET_SIZE_LIMIT) {
			return bucket
		}
	}
	return structs.Bucket{}
}