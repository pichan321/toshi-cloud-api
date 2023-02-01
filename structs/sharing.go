package structs

type ShareFile struct{
	Handle string `json:"handle"`
	Owner string `json:"owner"`
	Recipient string `json:"recipient"`
}

type SharedUser struct {
	Uuid string `json:"uuid"`
	Username string `json:"username"`
	Handle string `json:"handle"`
	IsShared string `json:"shared" db:"shared"`
}