package service

import (
	"helloworld/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateFood_Ingredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req model.Food_Ingredient

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "cannot input", "eror": err.Error()})
			return
		}
		// var findFood model.Food
		// var findingredient model.Ingredient
		// var findName model.Food_Ingredient
		// fmt.Print(req.FoodID)
		// if err := dbInstant.Model(&model.Food{}).Select("name").Where("id = ?", req.FoodID).Find(&findFood).Error; err != nil {
		// 	fmt.Print(err)
		// 	c.JSON(400, gin.H{"message": "failed to get food_name"})
		// 	return
		// }
		// if err := dbInstant.Model(&model.Ingredient{}).Select("name").Where("id = ?", req.IngredientID).Find(&findingredient).Error; err != nil {
		// 	c.JSON(400, gin.H{"message": "failed to get food_ingredient"})
		// 	return
		// }

		// findName = model.Food_Ingredient{
		// 	FoodID:         req.FoodID,
		// 	// FoodName:       findFood.Name,
		// 	IngredientID:   req.IngredientID,
		// 	// IngredientName: findingredient.Name,
		// }

		if err := dbInstant.Create(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to create food_ingredient", "error": err.Error()})
			return
		}

		c.JSON(200, req)
	}
}

func ViewFood_Ingredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		var res []model.Food_Ingredient

		if err := dbInstant.Model(&model.Food_Ingredient{}).Find(&res).Error; err != nil {
			c.JSON(400, gin.H{"message": "failed to get food_ingredient"})
			return
		}

		c.JSON(200, res)
	}
}

func UpdateFood_Ingredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req model.Food_Ingredient
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "cannot input", "eror": err.Error()})
			return
		}

		if err := dbInstant.Where("id = ?", dataID).Updates(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to get foods", "error": err.Error()})
			return
		}

		c.JSON(200, dataID)
	}
}

func DeleteFood_Ingredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req model.Food_Ingredient

		if err := dbInstant.Where("id = ?", dataID).Delete(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to get foods", "error": err.Error()})
			return
		}

		c.JSON(200, dataID)
	}
}
