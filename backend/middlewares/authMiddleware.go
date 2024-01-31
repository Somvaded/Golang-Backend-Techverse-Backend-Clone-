package middlewares

import (
	"backend/models"
	"backend/util"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Protect(c *gin.Context){
	jwtToken,err:= c.Cookie("jwt")
	if(err!=nil){
		c.JSON(http.StatusUnauthorized,err.Error())
		c.Abort()
		return
	}
	ok,err:=util.AuthorizeToken(jwtToken)
	if(ok){
		c.Next()
	}else{
		c.AbortWithError(http.StatusUnauthorized,err)
	}
}

func Admin(c *gin.Context){
	var user models.User
	err:= c.ShouldBindJSON(&user)
	if(err!=nil){
		c.AbortWithError(http.StatusBadGateway,err)
		return
	}
	if(user.IsAdmin){
		c.Next()
	}
	c.AbortWithError(http.StatusUnauthorized,errors.New("not admin"))

}