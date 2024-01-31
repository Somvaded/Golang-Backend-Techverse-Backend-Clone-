package main

import (
	"backend/controllers"
	// "backend/database"
	"backend/interfaces"
	"context"
	"fmt"
	"log"

	// "net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var(
	server *gin.Engine
	userMethods interfaces.UserMethods
	UserController controllers.UserController
	UserCollection *mongo.Collection
	Client *mongo.Client
	ctx context.Context
)
func init(){
	const uri = "mongodb+srv://sovajitr:Theanimeman@cluster0.yekwmys.mongodb.net/?retryWrites=true&w=majority"
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	Client, err := mongo.Connect(context.TODO(), opts)
	// Client.Ping(context.TODO(),nil)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("mongo conn established")
	ctx = context.TODO()
	server = gin.Default()
	UserCollection = Client.Database("TODO").Collection("users")
	userMethods=interfaces.UserMethodConst(UserCollection,ctx)
	UserController = controllers.New(userMethods)
}
func main(){
	defer Client.Disconnect(ctx)

	basepath:= server.Group("/api")
	UserController.RegisterUserRoutes(basepath)
	log.Fatal(server.Run(":8080"))	
}
