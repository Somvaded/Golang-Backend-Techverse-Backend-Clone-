package models

import "github.com/golang-jwt/jwt"

type CustomClaims struct{
	Id string
	isAdmin bool 
	jwt.StandardClaims
}