package models

import (
	"time"
)

type Review struct{
	UserId string `json:"userid" bson:"userid"`
	Username string `json:"username" bson:"username"`
	Comment string `json:"comment" bson:"comment"`
	Stars int `json:"stars" bson:"stars"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
} 
type Product struct{
	UserId      string             `json:"userid,omitempty" bson:"userid,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
    Image       string             `bson:"image,omitempty" json:"image,omitempty"`
    Brand       string             `bson:"brand,omitempty" json:"brand,omitempty"`
    Category    string             `bson:"category,omitempty" json:"category,omitempty"`
    Description string             `bson:"description,omitempty" json:"description,omitempty"`
    Price       float64            `bson:"price,omitempty" json:"price,omitempty"`
    Rating      int                `bson:"rating,omitempty" json:"rating,omitempty"`
    NumReviews  int                `bson:"numReviews,omitempty" json:"numReviews,omitempty"`
    CountInStock int               `bson:"countInStock,omitempty" json:"countInStock,omitempty"`
	Reviews    []Review            `bson:"reviews,omitempty" json:"reviews,omitempty"`
	CreatedAt  time.Time           `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time           `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type ProductResponse struct{
	UserId      string             `json:"userid,omitempty" bson:"userid,omitempty"`
	Name        string             `json:"name" bson:"name"`
    Image       string             `bson:"image" json:"image"`
    Brand       string             `bson:"brand" json:"brand"`
    Category    string             `bson:"category" json:"category"`
    Description string             `bson:"description" json:"description"`
    Price       float64            `bson:"price" json:"price"`
    Rating      int                `bson:"rating" json:"rating"`
    NumReviews  int                `bson:"numReviews" json:"numReviews"`
    CountInStock int               `bson:"countInStock" json:"countInStock"`
	Reviews    []Review            `bson:"reviews" json:"reviews"`
	CreatedAt  time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt" bson:"updatedAt"`
}



