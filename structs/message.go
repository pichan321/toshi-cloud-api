package structs

type Message struct {
	Message string `json:"message"`
	Code    int32  `json:"code"`
	Error   string `json:"error,omitempty"`
	Name    string `json:"name,omitempty"`
}