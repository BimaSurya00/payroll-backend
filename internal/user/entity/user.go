package entity

import (
	"github.com/itsahyarr/go-fiber-boilerplate/shared/entity"
)

type User struct {
	entity.BaseEntity `bson:",inline"`
	Name              string `json:"name" bson:"name"`
	Email             string `json:"email" bson:"email"`
	Password          string `json:"-" bson:"password"`
	Role              string `json:"role" bson:"role"`
	IsActive          bool   `json:"is_active" bson:"is_active"`
}

func NewUser(name, email, password, role string) *User {
	return &User{
		BaseEntity: entity.NewBaseEntity(),
		Name:       name,
		Email:      email,
		Password:   password,
		Role:       role,
		IsActive:   true,
	}
}
