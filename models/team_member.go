package models

import "time"

type TeamMember struct {
	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"-" bson:"password_hash"`
	Name         string    `json:"name" bson:"name"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}
