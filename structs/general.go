package structs

type File struct {
	Uuid         string  `json:"uuid" db:"uuid"`
	Name         string  `json:"name" db:"name"`
	Size         string  `json:"size" db:"size"`
	SizeMb       float64 `json:"sizeMb" db:"size_mb"`
	UploadedDate string  `json:"uploadedDate" db:"uploaded_date"`
	UserUuid     string  `json:"userUuid" db:"accounts_id"`
	BucketUuid   string  `json:"bucketUuid" db:"bucket_uuid"`
}

type Bucket struct {
	Uuid        string  `json:"uuid" db:"uuid"`
	Name        string  `json:"name" db:"name"`
	AccessToken string  `json:"accessToken" db:"access_token"`
	Size        float32 `json:"size" db:"size"`
}
