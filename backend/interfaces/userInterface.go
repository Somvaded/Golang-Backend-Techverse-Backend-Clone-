package interfaces

import (
	"backend/models"
	"backend/util"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserMethods interface {
	CreateUser(*models.User) (string,error)
	GetUser(*string) (*models.UserResponse,error)
	GetAllUsers() ([]*models.User,error)
	UpdateUser(*models.User,*string) (*models.UserResponse,error)
	DeleteUser(*string) error
	GetUserByEmail(*string,*string)(*models.User,error)
}

type UserMethodsImpl struct{
	userCollection *mongo.Collection
	ctx context.Context
}
func UserMethodConst(userCollection *mongo.Collection,ctx context.Context) UserMethods{
	return &UserMethodsImpl{
		userCollection: userCollection,
		ctx : ctx,
	}
} 


func (um *UserMethodsImpl) CreateUser(data *models.User) (string,error){
	var user models.User
	found:= um.userCollection.FindOne(um.ctx,bson.D{bson.E{Key: "email", Value: data.Email}}).Decode(&user)
	if found==mongo.ErrNoDocuments{
		result, err := um.userCollection.InsertOne(um.ctx,data)
		userId:= result.InsertedID.(primitive.ObjectID)
		userIdString:= userId.Hex()
		return userIdString,err
	}
	return "",errors.New("user already exists")
} 

func (um *UserMethodsImpl) GetUserByEmail(email *string, password *string)(*models.User,error){
	var user *models.User
	query := bson.D{bson.E{Key: "email" , Value: email}}
	err := um.userCollection.FindOne(um.ctx,query).Decode(&user)
	if((err!=nil)){
		return nil,err
	}
	newUser:= &models.User{
		Id: user.Id,
		UserName: user.UserName,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
		Password: user.Password,
	}
	return newUser,err
}

func (um *UserMethodsImpl) GetUser(id *string) (*models.UserResponse,error){
	var user *models.User
	userid,err:= primitive.ObjectIDFromHex(*id)
	if(err!=nil){
		return nil,err
	}
	query := bson.D{bson.E{Key: "_id" , Value: userid}}
	omit := bson.M{"password":0}
	err= um.userCollection.FindOne(um.ctx, query,options.FindOne().SetProjection(omit)).Decode(&user)
	response:= &models.UserResponse{
		Id: user.Id,
		UserName: user.UserName,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
	}
	return response,err
}

func (um *UserMethodsImpl) GetAllUsers() ([]*models.User,error){
	var data []*models.User
	omit := bson.M{"email":0}
	cursor,err:= um.userCollection.Find(um.ctx,bson.D{{}},options.Find().SetProjection(omit))
	if(err!=nil){
		return nil,err
	}
	for cursor.Next(um.ctx){
		var user models.User
        err := cursor.Decode(&user)
        if err != nil {
            return nil,err
        }

        data =append(data, &user)
	}
	if err:= cursor.Err(); err !=nil{
		return nil,err
	}
	cursor.Close(um.ctx)

	if(len(data)==0){
		return nil, mongo.ErrNoDocuments
	}
	return data,nil
}

func (um *UserMethodsImpl) UpdateUser(data *models.User,userid *string) (*models.UserResponse,error){
	var user models.User
	id,err:= primitive.ObjectIDFromHex(*userid)
	if(err!=nil){
		return nil,errors.New("not object id")
	}
	filter := bson.D{bson.E{Key: "_id" , Value: id}}
	err=um.userCollection.FindOne(um.ctx,filter).Decode(&user)
	if(err!=nil){
		return nil,err
	}
	update := bson.M{}
	response:=&models.UserResponse{}
	if(data.UserName!=""){
		update["username"]=data.UserName
		response.UserName=data.UserName
	} else{
		response.UserName=user.UserName
	}
	if(data.Email!=""){
		update["email"]=data.Email
		response.Email=data.Email
	}else{
		response.Email=user.Email
	}
	if(data.Password!=""){
		hashedId,err:=util.Encrypt(data.Password)
		if(err!=nil){
			return nil,err
		}
		update["password"]=hashedId
	}
	response.Id=user.Id
	response.IsAdmin=user.IsAdmin
	update["updated_at"]=time.Now()
	result,err:=um.userCollection.UpdateByID(um.ctx,id,bson.M{"$set":update})
	if(err!=nil){
		return nil,err
	}
	if(result.MatchedCount!=1){
		return nil,mongo.ErrNoDocuments
	}
	return response,nil
}

func (um *UserMethodsImpl) DeleteUser(id *string) (error){
	userid,err:= primitive.ObjectIDFromHex(*id)
	if(err!=nil){
		return errors.New("not object id")
	}
	filter := bson.D{bson.E{Key: "_id" , Value: userid}}
	result,err := um.userCollection.DeleteOne(um.ctx,filter)
	if(result.DeletedCount ==0){
		return mongo.ErrNoDocuments
	}
	return err
}


