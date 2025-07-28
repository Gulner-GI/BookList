package models

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type Update struct {
	Task *string `json:"task"`
	Done *bool   `json:"done"`
}
