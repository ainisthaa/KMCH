package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"helloworld/conf" 
	"helloworld/internal/dbconn"
	"helloworld/service"
)

func main() {

	confic, err := conf.NewConfig()
	if err != nil {
		return
	}
	dbInstant := dbconn.DBConnect(confic.SERVICE_DB_USER, confic.SERVICE_DB_PASS, confic.SERVICE_DB_HOST, confic.SERVICE_DB_PORT, confic.SERVICE_DB_NAME)

	// dbInstant.Migrator().DropTable(&model.Food_Ingredient{})
	// dbInstant.AutoMigrate(&model.Food_Ingredient{})

	server := gin.Default()

	// server.Use(cors.Default()) // อนุญาตทุก origin (*)
	// server.Use()
	fmt.Print("Hello")
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	server.POST("category", service.CreateCategory(dbInstant))
	server.POST("food", service.CreateFood(dbInstant))
	server.POST("ingredient", service.CreateIngredient(dbInstant))
	server.POST("foodingredient", service.CreateFood_Ingredient(dbInstant))
	server.GET("ingredient", service.ViewIngredient(dbInstant))
	server.GET("category", service.ViewCategory(dbInstant))
	server.GET("food", service.ViewFood(dbInstant))
	server.GET("foodingredient", service.ViewFood_Ingredient(dbInstant))
	server.GET("sumcat", service.ViewCategory_Number(dbInstant))
	server.GET("getmiso", service.FindMiso(dbInstant))
	server.GET("noneedegg", service.NoneedEgg(dbInstant))
	server.GET("ingredient/:food_id", service.GetIngredientByFoodID(dbInstant))
	server.GET("food/:food_id/ingredient", service.GetFoodDetail(dbInstant))
	server.PUT("category/:id", service.UpdateCategory(dbInstant))
	server.PUT("ingredient/:id", service.UpdateIngredient(dbInstant))
	server.PUT("food/:id", service.UpdateFood(dbInstant))
	server.DELETE("category/:id", service.DeleteCategory(dbInstant))
	server.DELETE("food/:id", service.DeleteFood(dbInstant))
	server.DELETE("ingredient/:id", service.DeleteIngredient(dbInstant))

	server.Run(":8890")
}
