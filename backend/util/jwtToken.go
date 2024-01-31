package util

import (
	"backend/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var secret = []byte("lmao")
func GenerateToken(User models.User,userId string) (string,error){
	claims := &models.CustomClaims{
		Id: User.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	tokenString,err:= token.SignedString(secret)
	if(err!=nil){
		return "",err
	}
	return tokenString,nil
}

func AuthorizeToken(jwtToken string) (bool,error){
	claims:= &models.CustomClaims{}
	token,err := jwt.ParseWithClaims(jwtToken,claims,func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if(err==jwt.ErrSignatureInvalid){
		return false,err
	}
	if(!token.Valid){
		return false,errors.New("Invalid Token")
	}
	return true,nil
}