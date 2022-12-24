package structs

type FileParse struct {
	Filename string `json:"name"`
	Content  string `json:"content"`
	Type     string `json:"type"`
	User     string `json:"user"`
}