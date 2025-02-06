package main

import (
	"backend/controllers"
	"backend/util"
	"backend/interfaces"
	"context"
	"log"
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
	
	
	
}
func main(){
	util.Load()
	client := connectDB()
	ctx = context.TODO()
	defer Client.Disconnect(ctx)
	server = gin.Default()
	basepath:= server.Group("/api")

	// Initializing
	UserCollection = client.Database("TODO").Collection("users")
	userMethods=interfaces.UserMethodConst(UserCollection,ctx)
	UserController = controllers.New(userMethods)
	UserController.RegisterUserRoutes(basepath)

	ProductCollection := client.Database("TODO").Collection("products")
	ProductMethods := interfaces.ProductMethodConx(ProductCollection,ctx)
	ProductContoller := controllers.ProductControllerContx(ProductMethods)
	ProductContoller.RegisterProductRoutes(basepath)
	log.Fatal(server.Run(":8080"))	
}


func connectDB() *mongo.Client{
	uri:=util.Mongo_uri
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	Client, err := mongo.Connect(context.TODO(), opts)
	// Client.Ping(context.TODO(),nil)
	if err != nil {
		log.Fatal(err)
	}
	return Client
}
