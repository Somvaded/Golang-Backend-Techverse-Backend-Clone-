package models

import "github.com/golang-jwt/jwt"

type CustomClaims struct{
	Id string
	IsAdmin bool 
	jwt.StandardClaims
}