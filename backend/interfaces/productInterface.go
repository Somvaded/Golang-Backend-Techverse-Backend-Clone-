package interfaces

import (
	"backend/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type ProductMethods interface {
	CreateProduct(string) (*models.ProductResponse, error)
	GetProduct(string) ([]*models.ProductResponse, error)
	GetProductById(string) (*models.ProductResponse, error)
	UpdateProduct(*models.Product, string) (*models.ProductResponse, error)
	GetTopProducts() ([]*models.ProductResponse, error)
	CreateProductReview(*models.Review,string) error
	DeleteProduct(string) error
}

type ProductMethodsImpl struct {
	ProductCollection *mongo.Collection
	ctx               context.Context
}

// CreateProductReview implements ProductMethods.
func (pm *ProductMethodsImpl) CreateProductReview(review *models.Review, productid string) error {
	id,err  := primitive.ObjectIDFromHex(productid)
	if err!= nil {
		return err
	}
	
	filter := bson.D{{Key: "_id",Value: id}}
	result := pm.ProductCollection.FindOne(pm.ctx,filter)
	if(result.Err() ==mongo.ErrNoDocuments){
		return errors.New("product does not exists")
	}
	var product models.Product
	err = result.Decode(&product)
	if(err != nil){
		return err
	}
	productReviews := product.Reviews
	for i := 0; i < len(productReviews); i++ {
		if(	productReviews[i].UserId == review.UserId){
			return errors.New("product review by this user already exists")
		}
	}
	product.Rating = ((product.Rating*product.NumReviews)+review.Stars)/(product.NumReviews+1)

	productReviews = append(productReviews, *review)
	product.NumReviews = len(productReviews)
	product.Reviews = productReviews

	filter1 := bson.M{"$set":product}

	_ , err = pm.ProductCollection.UpdateOne(pm.ctx , filter , filter1)

	if(err != nil){
		log.Println("here")
		return err
	}
	return nil
}

func ProductMethodConx(productCollection *mongo.Collection, ctx context.Context) ProductMethods {
	return &ProductMethodsImpl{
		ProductCollection: productCollection,
		ctx:               ctx,
	}
}

// GetProduct implements ProductMethods.
func (pm *ProductMethodsImpl) GetProduct(searchQuery string) ([]*models.ProductResponse, error) {
	var result []*models.ProductResponse

	cursor, err := pm.ProductCollection.Find(pm.ctx, bson.M{"name": bson.M{"$regex": searchQuery, "$options": "i"}})

	if err != nil {
		log.Println("Error finding the products:", err)
		return nil, err
	}
	err = cursor.All(pm.ctx, &result)

	if err != nil {
		log.Println("Error decoding the products:", err)
		return nil, err
	}
	return result, nil
}

// GetProductById implements ProductMethods.
func (pm *ProductMethodsImpl) GetProductById(id string) (*models.ProductResponse, error) {
	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("error converting ObjectID")
	}
	// Ensure MongoDB collection and context are initialized
	if pm.ProductCollection == nil {
		return nil, errors.New("MongoDB collection is not initialized")
	}
	if pm.ctx == nil {
		return nil, errors.New("MongoDB context is missing")
	}

	filter := bson.D{bson.E{Key: "_id", Value: productID}}
	result := pm.ProductCollection.FindOne(pm.ctx, filter)

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, errors.New("product not found")
		}
		return nil, result.Err()
	}

	var product models.ProductResponse // Fix nil pointer issue
	err = result.Decode(&product)
	if err != nil {
		return nil, errors.New("failed to decode product")
	}

	return &product, nil
}

// GetTopProducts implements ProductMethods.
func (pm *ProductMethodsImpl) GetTopProducts() ([]*models.ProductResponse, error) {
	findOptions := options.Find().SetSort(bson.D{{Key: "rating", Value: -1}})
	findOptions.SetLimit(3)
	resultCur, err := pm.ProductCollection.Find(pm.ctx, bson.M{}, findOptions)
	if err != nil {
		log.Println("Error finding the top products:", err)
		return nil, err
	}
	var products []*models.ProductResponse
	err = resultCur.All(pm.ctx, products)
	if err != nil {
		log.Println("Error decoding the top products:", err)
		return nil, err
	}
	return products, err
}

// UpdateProduct implements ProductMethods.
func (pm *ProductMethodsImpl) UpdateProduct(product *models.Product, id string) (*models.ProductResponse, error) {
	productId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error getting the objectID from Hex:", err)
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: productId}}

	product.UpdatedAt = time.Now()
	update := bson.M{"$set": product}

	up, err := pm.ProductCollection.UpdateOne(pm.ctx, filter, update)

	if err != nil {
		log.Println("Error updating the products:", err)
		return nil, err
	}
	if up.MatchedCount == 0 {
		log.Println(" product not found:", err)
		return nil, mongo.ErrNoDocuments
	}

	var updatedProduct *models.ProductResponse
	err = pm.ProductCollection.FindOne(pm.ctx, filter).Decode(&updatedProduct)

	if err != nil {
		log.Println("Error finding the product:", err)
		return nil, err
	}
	return updatedProduct, nil

}

func (pm *ProductMethodsImpl) CreateProduct(userId string) (*models.ProductResponse, error) {

	newProduct := models.ProductResponse{
		Name:         "Sample Name",
		Brand:        "sample brand",
		Category:     "sample category",
		Price:        0,
		CountInStock: 0,
		NumReviews:   0,
		Reviews:      nil,
		Description:  "samlpe desc",
		UserId:       userId,
		Image:        "/image/sample.jpg",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	addedProduct, err := pm.ProductCollection.InsertOne(pm.ctx, newProduct)

	if err != nil {
		log.Println("Error inserting the product:", err)
		return nil, err
	}

	var updatedProduct *models.ProductResponse
	err = pm.ProductCollection.FindOne(pm.ctx, bson.D{{Key: "_id", Value: addedProduct.InsertedID}}).Decode(&updatedProduct)

	if err != nil {
		log.Println("Error finding the product:", err)
		return nil, err
	}
	return updatedProduct, nil
}

func (pm *ProductMethodsImpl) DeleteProduct(userId string) error {
	ID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: ID}}

	deletedProduct, err := pm.ProductCollection.DeleteOne(pm.ctx, filter)

	if err != nil {
		return err
	}
	if deletedProduct.DeletedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}
