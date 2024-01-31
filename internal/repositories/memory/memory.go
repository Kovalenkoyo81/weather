package memory

import (
	"github.com/Kovalenkoyo81/weather/internal/models"
)

type Repository struct {
	users []models.User
}

func (r *Repository) AddUser(user models.User) {
	r.users = append(r.users, user)
}

func (r *Repository) FindUser(name string) bool {
	for _, u := range r.users {
		if u.Name == name {
			return true
		}
	}
	return false
}
