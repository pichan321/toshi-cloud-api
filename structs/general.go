package structs

import "context"

type File struct {
	Uuid         string  `json:"uuid" db:"uuid" gorm:"primaryKey"`
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
	Hidden       bool    `json:"hidden" db:"hidden"`
	SharedFile   bool    `json:"shared_file"`
}

type Bucket struct {
	Uuid        string  `json:"uuid" db:"uuid"`
	Name        string  `json:"name" db:"name"`
	AccessToken string  `json:"accessToken" db:"access_token"`
	Size        float32 `json:"size" db:"size"`
	ShareLink   string  `json:"shareLink" db:"sharelink"`
}

type FileContent struct {
	Content string `json:"content"`
}

type CustomClaims struct {
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	Nickname      string `json:"nickname"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
	Sub           string `json:"sub"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}