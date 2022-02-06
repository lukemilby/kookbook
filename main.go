// Recipes API
//
//This is a sample recipes API. You can find out more about
//the API at https://github.com/PacktPublishing/BuildingDistributed-Applications-in-Gin.
//
// Schemes: http
// Host: localhost:5000
// BasePath: /
// Version: 1.0.0
// Contact: Luke Milby <luke.milby@gmail.com> http://lukemilby.com
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/lukemilby/kookbook/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
)

// swagger:operation GET /recipes/search recipes searchRecipes
// Returns list of recipes that match the tag
// ---
// parameters:
// - name: tag
// in: query
// description: tag to search by
// required: true
// type: string
// produces:
// - application/json
// responses:
// '200':
// description: Successful operation
//func SearchRecipesHandler(c *gin.Context) {
//	tag := c.Query("tag")
//	listOfRecipes := make([]models.Recipe, 0)
//
//	for i := 0; i < len(recipes); i++ {
//		found := false
//		for _, t := range recipes[i].Tags {
//			if strings.EqualFold(t, tag) {
//				found = true
//			}
//		}
//		if found {
//			listOfRecipes = append(listOfRecipes,
//				recipes[i])
//		}
//	}
//
//	c.JSON(http.StatusOK, listOfRecipes)
//}


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != os.Getenv("X_API_KEY"){
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler

func init() {
	// load env file
	err := godotenv.Load()
	/*
	MONGO_URI | Connection and authentication | string | "mongodb://admin:password@10.10.1.121:27017/demo?authSource=admin"
	MONGO_DATABASE | Database name |string | demo
	X_API_KEY | API Key for middleware |string | adfes-wervgse-wf831
	JWT_SECRET | Token used to sign claim | string | whatyouwanthere
	*/

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx = context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "10.10.1.121:6379",
		Password: "",
		DB: 0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipeHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipesHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.Run(":5000")
}