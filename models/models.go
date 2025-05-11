package models

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var Todos = []Todo{
	{ID: 1, Task: "Learn Go", Done: false},
}
