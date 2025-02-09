package controllers

import (
	"backend/interfaces"
	"backend/middlewares"
	"backend/models"
	"log"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductController struct {
	ProductMethods interfaces.ProductMethods
}

func ProductControllerContx (productMethods interfaces.ProductMethods) (*ProductController) {
	return &ProductController{
		ProductMethods : productMethods,
	}
}

func (pc *ProductController) CreateNewProduct (ctx *gin.Context){
	userId,ok := ctx.Get("userId")
	log.Println(userId)
	if !ok {
		log.Println("userId attachment in middleware error")
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : "userId attachment in middleware error",
		})
		return
	}
	product ,err := pc.ProductMethods.CreateProduct(userId.(string))
	if err!=nil {
		log.Println("createProduct interface error")
		ctx.JSON(http.StatusInternalServerError , gin.H{
			"message" : err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK , product)
}

func (pc *ProductController) GetProduct(ctx *gin.Context){
	searchQuery := ctx.Query("search")

	var products []*models.ProductResponse
	products , err := pc.ProductMethods.GetProduct(searchQuery)

	if err != nil {
		log.Fatal("getProduct error")
		ctx.JSON(http.StatusInternalServerError , gin.H{
			"message" : err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK , products)
} 

func (pc *ProductController) GetProductById (ctx *gin.Context){
	productId := ctx.Param("id")
	_ , err := primitive.ObjectIDFromHex(productId)
	if err!=nil {
        ctx.JSON(http.StatusOK, gin.H{"message": "Invalid ObjectID format"})
        return
    }

    product, err := pc.ProductMethods.GetProductById(productId)
    if err != nil {
        log.Println("GetProductById error:", err)

        status := http.StatusInternalServerError
        if err.Error() == "product not found" {
            status = http.StatusNotFound
        } else if strings.Contains(err.Error(), "invalid ObjectID format") {
            status = http.StatusBadRequest
        }

        ctx.JSON(status, gin.H{"message": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Product retrieved successfully",
        "product": product,
    })
}


func(pc *ProductController) UpdateProduct(ctx *gin.Context){
	productId := ctx.Param("id")
	_ , err := primitive.ObjectIDFromHex(productId)
	if err!=nil {
        ctx.JSON(http.StatusOK, gin.H{"message": "Invalid ObjectID format"})
        return
    }
	var productData models.Product 
	err = ctx.ShouldBindJSON(&productData)
	if( err != nil){
		log.Println("error decoding json" , err)
		ctx.JSON(http.StatusInternalServerError , gin.H{
			"message" : err.Error(),
		})
	}
	updatedProduct,err := pc.ProductMethods.UpdateProduct(&productData,productId)
	if err != nil {
		log.Println("error updating product" , err)
		ctx.JSON(http.StatusInternalServerError , gin.H{
			"message" : err.Error(),
		})
	}
	ctx.JSON(http.StatusOK , updatedProduct)
}


func (pc *ProductController) RegisterProductRoutes (rg *gin.RouterGroup){
	mainRoute := rg.Group("/products")
	//public 
	mainRoute.GET("/", pc.GetProduct)
	//admin
	adminRoute := mainRoute.Group("/admin" , middlewares.Protect, middlewares.Admin )
	adminRoute.POST("/create" , pc.CreateNewProduct)
	adminRoute.POST("/:id" , pc.UpdateProduct)
	adminRoute.GET("/:id", pc.GetProductById)
}