package structs

type Account struct {
	Uuid     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Token    string `json:"token"`
}