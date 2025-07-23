package model

type BookRes struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

type ChangeData struct {
	Title       *string `json:"title,omitempty"`
	Author      *string `json:"author,omitempty"`
	Description *string `json:"description,omitempty"`
}
