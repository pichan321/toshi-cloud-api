package utils

import (
	"file-api/cloud"
	"file-api/structs"
	"fmt"
)

const BUCKET_SIZE_LIMIT = 150000.00

func GetBucketUuid(fileSizeMb float64) structs.Bucket {
	db, _ := cloud.GetPostgres()
	defer db.Close()
	
	rows, _ := db.Queryx("select * from buckets")
	var bucket structs.Bucket
	for rows.Next(){

		rows.StructScan(&bucket)
		if (bucket.Size + float32(fileSizeMb) <= float32(BUCKET_SIZE_LIMIT)) {
			fmt.Println("COMPARE")
			fmt.Println((bucket.Size + float32(fileSizeMb)) <= float32(BUCKET_SIZE_LIMIT))
			fmt.Println("SIZE ABOUT TO UPLOAD")
			fmt.Println(bucket.Size + float32(fileSizeMb))
			fmt.Println("BUCKET CURRENT SIZE")
			fmt.Println(BUCKET_SIZE_LIMIT)
			fmt.Println("GET BUCKET")
			fmt.Printf("%v", bucket)

			break
		
		}
		fmt.Println("still finding")
	}
	return bucket
}