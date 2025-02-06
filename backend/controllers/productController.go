package controllers

import (
	"backend/interfaces"
	"backend/middlewares"
	"backend/models"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
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
	var user models.User 
	err :=ctx.ShouldBindJSON(&user)
	if err !=nil {
		log.Println("binding error")
		ctx.JSON(http.StatusBadRequest , gin.H{
			"message" : err.Error(),
		})
		return
	}
	product ,err := pc.ProductMethods.CreateProduct(&user)
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

	var products []*models.Product
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
	productId := ctx.Query("id")
	
	product , err := pc.ProductMethods.GetProductById(productId)
	if(err != nil ){
		log.Fatal("getproductbyid error")
		ctx.JSON(http.StatusNotFound ,gin.H{
			"message" : err.Error(),
		})
		return 
	}
	ctx.JSON(http.StatusOK , product)
}


func(pc *ProductController) UpdateProduct (ctx *gin.Context){
	var productData *models.Product 
	err := ctx.ShouldBindJSON(&productData)

	productID := ctx.Query("id")
	if( err != nil){
		log.Fatal("error decoding json" , err)
		ctx.JSON(http.StatusInternalServerError , gin.H{
			"message" : err.Error(),
		})
	}
	updatedProduct,err := pc.ProductMethods.UpdateProduct(productData,productID)
	if err != nil {
		log.Fatal("error updating product" , err)
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
	
	adminRoute := mainRoute.Group("/admin" , middlewares.Admin , middlewares.Protect)

	adminRoute.POST("/create" , pc.CreateNewProduct)
	adminRoute.POST("/" , pc.UpdateProduct)
	adminRoute.GET("/", pc.GetProductById)
}