// internal/models/businessModels.go

package models

type User struct {
	Name string `json:"name"`
}

type Favorite struct {
	City       string   `json:"city"`
	Parameters []string `json:"parameters"`
}

type LoginRequest struct {
	User string `json:"user"`
}
