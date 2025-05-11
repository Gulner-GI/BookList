package main

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var todos = []Todo{
	{ID: 1, Task: "Learn Go", Done: false},
}
