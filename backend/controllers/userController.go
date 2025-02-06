package controllers

import (
	"backend/interfaces"
	"backend/middlewares"
	"backend/models"
	"backend/util"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

type UserController struct{
	UserMethods interfaces.UserMethods
}
func New(UserMethods interfaces.UserMethods) UserController{
	return UserController{
		UserMethods: UserMethods,
	}
}

// PUBLIC ROUTES

func (uc *UserController) RegisterUser(ctx *gin.Context){
	var user models.User
	err:=ctx.ShouldBindJSON(&user);if err!= nil{
		log.Println("binding error")
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : err.Error(),
		})
		return
	}
	hashedPassword,err:= util.Encrypt(user.Password)
	fmt.Println(hashedPassword)
	if(err!=nil){
		log.Println("hashing error")
		ctx.JSON(http.StatusBadGateway,gin.H{
			"message": err.Error(),
		})
		return
	}
	user.Password=hashedPassword
	user.CreatedAt=time.Now()
	user.UpdatedAt=time.Now()
	userID,err := uc.UserMethods.CreateUser(&user)
	if(err!=nil){
		log.Println("mongo error")
		ctx.JSON(http.StatusBadGateway,gin.H{
			"message": err.Error(),
		})
		return
	}
	jwtToken,err:= util.GenerateToken(user,userID)
	if(err!=nil){
		ctx.JSON(http.StatusBadGateway,gin.H{
			"message": err.Error(),
		})
		return
	}
	response:= models.UserResponse{
		Id: userID,
		UserName: user.UserName,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
	}
	ctx.SetCookie("jwt",jwtToken,int(time.Hour*24),"/","localhost",false,true)
    ctx.JSON(http.StatusOK,response)
} 

func (uc *UserController) LoginUser(ctx *gin.Context){
	var user models.User
	err:=ctx.ShouldBindJSON(&user);if err!= nil{
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : err.Error(),
		})
		return
	}
	findUser,err:=uc.UserMethods.GetUserByEmail(&user.Email,&user.Password)
	if err!= nil{
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : err.Error(),
		})
		return
	}
	fmt.Println(findUser.Password)
	pass:=util.CheckPassword(findUser.Password,user.Password)
	if(!pass){
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : "wrong password",
		})
		return
	}
	jwtToken,err:= util.GenerateToken(*findUser,findUser.Id)
	if(err!=nil){
		fmt.Println("here")
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : err.Error(),
		})
		return
	}
	response:= models.UserResponse{
		Id : findUser.Id,
		UserName: findUser.UserName,
		Email: findUser.Email,
		IsAdmin: findUser.IsAdmin,
	}
	ctx.SetCookie("jwt",jwtToken,int(time.Hour*24),"/","localhost",false,true)
	ctx.JSON(http.StatusOK , response)
}


func(uc *UserController) LogoutUser (c *gin.Context){
	c.SetCookie("jwt","",0,"/","localhost",false,true)
	c.JSON(http.StatusOK,gin.H{"message":"User logged out Succesfully"})
}



//  PRIVATE ROUTES


func (uc *UserController) GetUserProfile(ctx *gin.Context){
	var User models.User
	err := ctx.ShouldBindJSON(&User)
	if(err!=nil) {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message":err.Error(),
		})
	}
	user ,err := uc.UserMethods.GetUser(&User.Id)
	if(err!=nil) {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message":err.Error(),
		})
	}
	response:= models.UserResponse{
		Id : user.Id,
		UserName: user.UserName,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
	}
	ctx.JSON(http.StatusOK,response)
}

func (uc *UserController) UpdateUserProfile(ctx *gin.Context){
	var user models.User
	err:=ctx.ShouldBindJSON(&user);if err!= nil{
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : err.Error(),
		})
		return
	}
	UpdatedUser,err:= uc.UserMethods.UpdateUser(&user,&user.Id)
	if(err!=nil){
		ctx.JSON(http.StatusBadGateway,gin.H{
			"message": err.Error(),
		})
		return
	}
    ctx.JSON(http.StatusOK,UpdatedUser)
}


// ADMIN ROUTES

func (uc *UserController) GetUser (ctx *gin.Context){
	userId := ctx.Param("id")
	user ,err := uc.UserMethods.GetUser(&userId)
	if(err!=nil) {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message":err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK,user)
}




func (uc *UserController) GetAllUsers(ctx *gin.Context) {
	data,err := uc.UserMethods.GetAllUsers()
	if(err!=nil){
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message":err.Error(),
		})
	}
	ctx.JSON(http.StatusOK,data)
}





func (uc *UserController) DeleteUser(ctx *gin.Context){
	var data models.User
	ctx.ShouldBindJSON(&data)
	userid:=ctx.Param("id")
	err:=uc.UserMethods.DeleteUser(&userid)
	if(err!=nil){
		ctx.JSON(http.StatusBadGateway,gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"message":"success",
	})
}


func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup){
	mainRoute := rg.Group("/users")

	// PUBLIC
	publicRoute := mainRoute
	publicRoute.POST("/",uc.RegisterUser)
	publicRoute.POST(("/login"),uc.LoginUser)

	// private
	userRoute:= mainRoute.Group("/",middlewares.Protect)
	
	userRoute.POST("/logout",uc.LogoutUser)
	userRoute.GET("/profile",uc.GetUserProfile)
	userRoute.PATCH("/profile",uc.UpdateUserProfile)
	

	// ADMIN
	adminRoute := mainRoute.Group("/",middlewares.Protect,middlewares.Admin)
	adminRoute.GET("/",uc.GetAllUsers)
	adminRoute.GET("/:id",uc.GetUser)
	adminRoute.DELETE("/:id",uc.DeleteUser)
}






