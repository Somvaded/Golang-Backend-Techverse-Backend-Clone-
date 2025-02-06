package interfaces

import (
	"backend/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductMethods interface {
	CreateProduct(*models.User) (*models.Product, error)
	GetProduct(string) ([]*models.Product, error)
	GetProductById(string) (*models.Product, error)
	UpdateProduct(*models.Product, string) (*models.Product, error)
	GetTopProducts() ([]*models.Product, error)
}

type ProductMethodsImpl struct {
	ProductCollection *mongo.Collection
	ctx               context.Context
}

func ProductMethodConx(productCollection *mongo.Collection, ctx context.Context) ProductMethods {
	return &ProductMethodsImpl{
		ProductCollection: productCollection,
		ctx:               ctx,
	}
}

// GetProduct implements ProductMethods.
func (pm *ProductMethodsImpl) GetProduct(searchQuery string) ([]*models.Product, error) {
	var result []*models.Product

	cursor, err := pm.ProductCollection.Find(pm.ctx , bson.M{ "name" : bson.M{"$regex" : searchQuery , "$options" : "i"}})

	if err != nil {
		log.Println("Error finding the products:" , err)
		return nil,err
	}
	err = cursor.All(pm.ctx , &result)

	if(err !=nil) {
		log.Println("Error decoding the products:" , err)
		return nil , err
	}
	return result , nil
}  

// GetProductById implements ProductMethods.
func (pm *ProductMethodsImpl) GetProductById(id string) (*models.Product, error) {
	productID := id
	filter := bson.M{"_id": productID}

	result := pm.ProductCollection.FindOne(pm.ctx , filter)
	if result.Err() != nil{
		log.Println("Error finding the product:" ,result.Err())
		return nil ,result.Err()
	}
	var product *models.Product
	err := result.Decode(product)
	if err != nil {
		log.Println("Error decoding the product:" , err)
		return nil ,err
	}
	return product , nil
}

// GetTopProducts implements ProductMethods.
func (pm *ProductMethodsImpl) GetTopProducts() ([]*models.Product, error) {
	findOptions := options.Find().SetSort(bson.D{{"rating",-1}})
	findOptions.SetLimit(3)
	resultCur ,err := pm.ProductCollection.Find(pm.ctx , bson.M{} , findOptions)
	if(err != nil){
		log.Println("Error finding the top products:" , err)
		return nil ,err
	}
	var products []*models.Product
	err = resultCur.All(pm.ctx , products)
	if err !=nil {
		log.Println("Error decoding the top products:" , err)
		return nil , err
	}
	return products , err
}

// UpdateProduct implements ProductMethods.
func (pm *ProductMethodsImpl) UpdateProduct(product *models.Product , id string ) (*models.Product, error) {
	updateData := &models.Product{
		Name:        product.Name,
		Image      : product.Image,  
		Brand      : product.Brand,
		Category    : product.Category,
		Description  : product.Description,
		Price   : product.Price,
		Rating    : product.Rating,
		NumReviews :product.NumReviews,
		CountInStock      : product.CountInStock,
		Reviews     : product.Reviews,
	}
	productId ,err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error getting the objectID from Hex:" , err)
		return nil ,err
	}
	filter := bson.M{"_id" : productId}
	update := bson.M{"$set" : updateData}

	up , err := pm.ProductCollection.UpdateByID(pm.ctx , filter ,update )

	if err!=nil {
		log.Println("Error updating the products:" , err)
		return nil ,err
	}
	if(up.MatchedCount == 0){
		log.Println(" product not found:" , err)
		return nil , mongo.ErrNoDocuments
	}

	var updatedProduct *models.Product
	err = pm.ProductCollection.FindOne(pm.ctx ,filter).Decode(&updatedProduct)

	if err!= nil {
		log.Println("Error finding the product:" , err)
		return nil ,err
	}
	return updatedProduct, nil

}



func (pm *ProductMethodsImpl) CreateProduct(user *models.User) (*models.Product, error) {

	newProduct := models.Product{
		Name:         "Sample Name",
		Brand:        "sample brand",
		Category:     "sample category",
		Price:        0,
		CountInStock: 0,
		NumReviews:   0,
		Description:  "samlpe desc",
		UserId:       user.Id,
		Image:        "/image/sample.jpg",
	}

	addedProduct, err := pm.ProductCollection.InsertOne(pm.ctx, newProduct)
	if err != nil {
		log.Println("Error inserting the product:" , err)
		return &newProduct, err
	}

	newProduct.ProductId = addedProduct.InsertedID.(primitive.ObjectID).Hex()
	return &newProduct, nil
}
