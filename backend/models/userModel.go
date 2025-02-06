package models

import (
	"time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	Id string `json:"id,omitempty" bson:"_id,omitempty"`
	UserName string `json:"username" bson:"username"`
	Email string    `json:"email" bson:"email"`
	IsAdmin bool    `json:"isAdmin" bson:"isAdmin"`
	Password string `json:"password," bson:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" bson:"updated_at"` 
}

type UserResponse struct{
	Id string `json:"id,omitempty" bson:"_id,omitempty"`
	UserName string `json:"username" bson:"username"`
	Email string    `json:"email" bson:"email"`
	IsAdmin bool    `json:"isAdmin" bson:"isAdmin"`	
}