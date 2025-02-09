package middlewares

import (
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
	ok,id,err:=util.AuthorizeToken(jwtToken)
	if(ok && err==nil){
		c.Set("userId",id)
		c.Next()
	}else {
        // If token is invalid, abort with an unauthorized error
        if err != nil {
            c.AbortWithError(http.StatusUnauthorized, err)
        } else {
            c.JSON(http.StatusUnauthorized, "Unauthorized access")
        }
    }
}

func Admin(c *gin.Context){
	jwtToken,err:= c.Cookie("jwt")
	if(err!=nil){
		c.JSON(http.StatusUnauthorized,err.Error())
		c.Abort()
		return
	}
	isAdmin , err := util.CheckAdmin(jwtToken)
	if(err!=nil){
		c.AbortWithError(http.StatusBadGateway,err)
		return
	}
	if(isAdmin){
		c.Next()
	} else if(!isAdmin) {
	    c.AbortWithError(http.StatusUnauthorized,errors.New("not Admin"))
	} else{
		c.AbortWithError(http.StatusUnauthorized,errors.New("error in Admin middleware"))
	}

}