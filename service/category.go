package service

import (
	"helloworld/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ForCount struct {
	CategoryID   uint   `json:"cat_id" gorm:"column:cat_id"`
	CategoryName string `json:"cat_name" gorm:"column:cat_name"`
	FoodCount    int64  `json:"food_count" gorm:"column:food_count"`
}

func CreateCategory(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req model.Category

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "cannot input", "error": err.Error()})
			return
		}

		if err := dbInstant.Create(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to create category", "error": err.Error()})
			return
		}

		c.JSON(200, req)
	}
}

func ViewCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err1 := strconv.Atoi(pageStr)
		limit, err2 := strconv.Atoi(limitStr)
		if err1 != nil || err2 != nil || page < 1 || limit < 1 {
			c.JSON(400, gin.H{"error": "invalid page or limit"})
			return
		}
		offset := (page - 1) * limit

		var categories []model.Category
		if err := db.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch categories"})
			return
		}

		var total int64
		if err := db.Model(&model.Category{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to count categories"})
			return
		}

		c.JSON(200, gin.H{
			"data":  categories,
			"total": total,
		})
	}
}

func ViewCategory_Number(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var cfood []ForCount

		if err := dbInstant.Table("categories").
			Select("categories.id as cat_id, categories.name as cat_name, count(foods.id) as food_count").
			Joins("LEFT JOIN foods ON categories.id = foods.category_id").
			Group("categories.id").
			Scan(&cfood).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch food count"})
			return
		}

		c.JSON(200, cfood)
	}
}

func UpdateCategory(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req model.Category
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "cannot input", "error": err.Error()})
			return
		}

		if err := dbInstant.Where("id = ?", dataID).Updates(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to update category", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "updated", "id": dataID})
	}
}

func DeleteCategory(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req model.Category
		if err := dbInstant.Where("id = ?", dataID).Delete(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to delete category", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "deleted", "id": dataID})
	}
}
