package structs

type Account struct {
	Uuid     string `json:"uuid,omitempty" gorm:"primaryKey"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Auth0 string `json:"auth0"`
	ApiKey string `json:"api_key" db:"api_key"`
}
