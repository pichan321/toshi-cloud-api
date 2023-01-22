package structs

type Account struct {
	Uuid     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Token    string `json:"token"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Token string `json:"token"` 
}