package models

type User struct {
	Name string `json:"name"`
}

type Favorite struct {
	City       string   `json:"city"`
	Parameters []string `json:"parameters"`
}
