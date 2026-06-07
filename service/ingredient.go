package service

import (
	"fmt"
	"helloworld/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateIngredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req model.Ingredient

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "cannot input", "eror": err.Error()})
			return
		}

		if err := dbInstant.Create(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to create food", "error": err.Error()})
			return
		}

		c.JSON(200, req)
	}
}

func ViewIngredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		// รับค่าจาก query string
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err1 := strconv.Atoi(pageStr)
		limit, err2 := strconv.Atoi(limitStr)
		if err1 != nil || err2 != nil || page < 1 || limit < 1 {
			c.JSON(400, gin.H{"error": "invalid page or limit"})
			return
		}

		offset := (page - 1) * limit

		var res []model.Ingredient
		if err := dbInstant.
			Limit(limit).
			Offset(offset).
			Find(&res).Error; err != nil {
			c.JSON(400, gin.H{"message": "failed to get ingredient", "error": err.Error()})
			return
		}

		var total int64
		if err := dbInstant.Model(&model.Ingredient{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{"message": "failed to count ingredients", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"data":  res,
			"total": total,
		})
	}
}

func UpdateIngredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req model.Ingredient
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

func DeleteIngredient(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req model.Ingredient

		if err := dbInstant.Where("id = ?", dataID).Delete(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to get foods", "error": err.Error()})
			return
		}

		c.JSON(200, dataID)
	}
}

type FoodIngredientDetail struct {
	IngredientID   int    `json:"ingredient_id"`
	IngredientName string `json:"ingredient_name"`
}

func GetIngredientByFoodID(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		foodID := c.Param("food_id")

		var res []FoodIngredientDetail
		fi_res := dbInstant.Model(&model.Food_Ingredient{}).Select("ingredients.id as ingredient_id, ingredients.name as ingredient_name").
			Where("food_id = ?", foodID).
			Joins("JOIN ingredients ON food_ingredients.ingredient_id = ingredients.id").
			Find(&res)
		if fi_res.Error != nil {
			c.JSON(400, gin.H{"message": "Failed to get food ingredients", "error": fi_res.Error.Error()})
			return
		}
		fmt.Println(res)
		c.JSON(200, res)
	}
}
