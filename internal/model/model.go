package model

type BookRes struct {
	Title       string `json:"title" required:"true"`
	Author      string `json:"author" required:"true"`
	Description string `json:"description" required:"true"`
}

type ChangeData struct {
	Title       *string `json:"title,omitempty"`
	Author      *string `json:"author,omitempty"`
	Description *string `json:"description,omitempty"`
}
