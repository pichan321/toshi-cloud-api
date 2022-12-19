package structs

type File struct {
	Uuid         string  `json:"uuid" db:"uuid"`
	Name         string  `json:"name" db:"name"`
	Size         string  `json:"size" db:"size"`
	SizeMb       float64 `json:"sizeMb" db:"size_mb"`
	UploadedDate string  `json:"uploadedDate" db:"uploaded_date"`
	UserUuid     string  `json:"userUuid" db:"account_uuid"`
	BucketUuid   string  `json:"bucketUuid" db:"bucket_uuid"`
	Part         int64   `json:"part"`
	Total        int64   `json:"total"`
	Status       string  `json:"status"`
	UploadID     string  `json:"uploadId"`
}

type Bucket struct {
	Uuid        string  `json:"uuid" db:"uuid"`
	Name        string  `json:"name" db:"name"`
	AccessToken string  `json:"accessToken" db:"access_token"`
	Size        float32 `json:"size" db:"size"`
	ShareLink   string  `json:"shareLink" db:"sharelink"`
}
