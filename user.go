package main

import (
	"time"

	"github.com/Isudin/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	statusCode int
	Id         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (u *User) getStatusCode() int {
	return u.statusCode
}

func mapUser(dbUser database.User) User {
	return User{
		Id:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}
