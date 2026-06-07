package service

import (
	"helloworld/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateFoodRequestBody struct {
	Name          string `json:"name"`
	Price         uint   `json:"price"`
	CategoryID    int    `json:"category_id"`
	IngredientsID []int  `json:"ingredients_id"`
}
type UpdatesFoodRequest struct {
	Name          string `json:"name"`
	Price         uint   `json:"price"`
	CategoryID    int    `json:"category_id"`
	IngredientsID []int  `json:"ingredients_id"`
}
type FoodDetail struct {
	Name        string                 `json:"name"`
	Price       float64                `json:"price"`
	CategoryID  int                    `json:"category_id"`
	Ingredients []FoodIngredientDetail `json:"ingredients"`
}

func CreateFood(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req CreateFoodRequestBody

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "cannot input", "error": err.Error()})
			return
		}
		var food_create model.Food = model.Food{
			Name:       req.Name,
			Price:      req.Price,
			CategoryID: req.CategoryID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		food_create_res := dbInstant.Create(&food_create)
		if food_create_res.Error != nil {
			c.JSON(400, gin.H{"message": "Failed to create food", "error": food_create_res.Error.Error()})
			return
		}

		var fi_create []model.Food_Ingredient = []model.Food_Ingredient{}
		// req.IngredientsID = [1, 2, 3]
		// ingredientID = 1 => 2 => 3
		// fi_create = [
		// 	{
		// 		FoodID:       19,
		//  	IngredientID: 1,
		// 		CreatedAt:    time.Now(),
		// 	},
		// 	{
		// 		FoodID:       19,
		//  	IngredientID: 2,
		// 		CreatedAt:   time.Now(),
		// 	},
		// ]
		for _, ingredientID := range req.IngredientsID {
			fi_create = append(fi_create, model.Food_Ingredient{
				FoodID:       food_create.ID,
				IngredientID: ingredientID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			})
		}
		food_ingredient_create_res := dbInstant.Create(&fi_create)
		if food_ingredient_create_res.Error != nil {
			c.JSON(400, gin.H{"message": "Failed to create food ingredients", "error": food_ingredient_create_res.Error.Error()})
			return
		}

		c.JSON(200, req)
	}
}

type ViewFoodResponse struct {
	Name        string                 `json:"name"`
	Price       float64                `json:"price"`
	CategoryID  int                    `json:"category_id"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Ingredients []FoodIngredientDetail `json:"ingredients"`
	ID          int                    `json:"id"`
}

func ViewFood(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ---------- 1. ดึง query params ----------
		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")
		q := strings.ToLower(strings.TrimSpace(c.Query("q")))
		minStr := c.Query("minPrice")
		maxStr := c.Query("maxPrice")
		sortKey := c.DefaultQuery("sort", "")
		var orderExpr string // ข้อความ SQL สำหรับ ORDER BY

		// ---------- 2. ดึง list ของ ingredient ที่ผู้ใช้เลือก ----------
		ingList := make([]string, 0)
		for _, v := range c.QueryArray("ingredient") {
			v = strings.ToLower(strings.TrimSpace(v))
			if v != "" {
				ingList = append(ingList, v)
			}
		}

		// ---------- 3. ตรวจสอบ page / limit ----------
		page, err1 := strconv.Atoi(pageStr)
		limit, err2 := strconv.Atoi(limitStr)
		if err1 != nil || err2 != nil || page < 1 || limit < 1 {
			c.JSON(400, gin.H{"error": "invalid page or limit"})
			return
		}
		offset := (page - 1) * limit

		// ---------- 4. เริ่มประกอบ query ----------
		base := db.Model(&model.Food{}).Where("foods.deleted_at IS NULL")

		if q != "" {
			base = base.Where("LOWER(foods.name) LIKE ?", "%"+q+"%")
		}

		// ---------- 5. เงื่อนไขวัตถุดิบ ----------
		if len(ingList) > 0 {
			base = base.Joins(`
				JOIN food_ingredients fi ON fi.food_id = foods.id AND fi.deleted_at IS NULL
				JOIN ingredients i       ON i.id       = fi.ingredient_id AND i.deleted_at IS NULL`).
				Where("LOWER(i.name) IN ?", ingList).
				Group("foods.id").
				Having("COUNT(DISTINCT i.id) = ?", len(ingList))
		}

		// ---------- 6. เงื่อนไขช่วงราคา ----------
		if minStr != "" || maxStr != "" {
			var conditions []string
			var values []interface{}

			if minStr != "" {
				if minVal, err := strconv.Atoi(minStr); err == nil {
					conditions = append(conditions, "foods.price >= ?")
					values = append(values, minVal)
				}
			}
			if maxStr != "" {
				if maxVal, err := strconv.Atoi(maxStr); err == nil {
					conditions = append(conditions, "foods.price <= ?")
					values = append(values, maxVal)
				}
			}
			if len(conditions) > 0 {
				base = base.Where(strings.Join(conditions, " AND "), values...)
			}
		}

		switch sortKey {
		case "name":
			orderExpr = "LOWER(foods.name) ASC"
		case "-name":
			orderExpr = "LOWER(foods.name) DESC"
		case "price":
			orderExpr = "foods.price ASC"
		case "-price":
			orderExpr = "foods.price DESC"
		case "category":
			orderExpr = "foods.category_id ASC"
		case "-category":
			orderExpr = "foods.category_id DESC"
		case "updated":
			orderExpr = "foods.updated_at ASC"
		case "-updated":
			orderExpr = "foods.updated_at DESC"
		default:
			orderExpr = "foods.id ASC"
		}

		// ---------- 7. นับจำนวนทั้งหมด ----------
		var total int64
		if err := base.Select("foods.id").Count(&total).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// ---------- 8. ดึงข้อมูลตามหน้า ----------
		var foods []model.Food

		if err := base.Select("foods.*").
			Order(orderExpr).
			Limit(limit).Offset(offset).
			Find(&foods).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// ---------- 9. รวม ingredients ของแต่ละอาหาร ----------
		var resp []ViewFoodResponse
		for _, f := range foods {
			var ings []FoodIngredientDetail
			if err := db.Table("food_ingredients fi").
				Select("i.id AS ingredient_id, i.name AS ingredient_name").
				Joins("JOIN ingredients i ON i.id = fi.ingredient_id AND i.deleted_at IS NULL").
				Where("fi.food_id = ? AND fi.deleted_at IS NULL", f.ID).
				Scan(&ings).Error; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			resp = append(resp, ViewFoodResponse{
				ID:          f.ID,
				Name:        f.Name,
				Price:       float64(f.Price),
				CategoryID:  f.CategoryID,
				UpdatedAt:   f.UpdatedAt,
				Ingredients: ings,
			})
		}

		// ---------- 10. ส่ง response ----------
		c.JSON(200, gin.H{
			"data":  resp,
			"total": total,
		})
	}
}

func GetFoodDetail(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		foodID := c.Param("food_id")

		var req model.Food
		if err := db.First(&req, foodID).Error; err != nil {
			c.JSON(404, gin.H{"error": "ไม่พบอาหาร"})
			return
		}

		// ดึงวัตถุดิบ
		var ingredients []FoodIngredientDetail
		db.Model(&model.Food_Ingredient{}).
			Select("ingredients.id as ingredient_id, ingredients.name as ingredient_name").
			Where("food_id = ?", foodID).
			Joins("JOIN ingredients ON food_ingredients.ingredient_id = ingredients.id").
			Find(&ingredients)

		detail := FoodDetail{
			Name:        req.Name,
			Price:       float64(req.Price),
			CategoryID:  req.CategoryID,
			Ingredients: ingredients,
		}

		c.JSON(200, detail)
	}
}

type ForChoose struct {
	FoodName string `json:"food_name"`
}

func NoneedEgg(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "3")

		page, err1 := strconv.Atoi(pageStr)
		limit, err2 := strconv.Atoi(limitStr)

		if err1 != nil || err2 != nil || page < 1 || limit < 1 {
			c.JSON(400, gin.H{"error": "invalid page or limit"})
			return
		}

		offset := (page - 1) * limit

		var res []ForChoose

		if err := dbInstant.Table("foods").
			Select("foods.name as food_name").
			Joins("JOIN food_ingredients ON food_ingredients.food_id = foods.id").
			Joins("JOIN ingredients ON ingredients.id = food_ingredients.ingredient_id").
			Group("foods.id").
			Having("SUM(CASE WHEN ingredients.name = ? THEN 1 ELSE 0 END) = 0", "egg").
			Limit(limit).
			Offset(offset).
			Scan(&res).Error; err != nil {

			c.JSON(500, gin.H{"error": "Failed to fetch food count"})
			return

		}
		c.JSON(200, gin.H{
			"page":  page,
			"limit": limit,
			"data":  res,
		})
	}
}

func FindMiso(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {

		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "3")

		page, err1 := strconv.Atoi(pageStr)
		limit, err2 := strconv.Atoi(limitStr)

		if err1 != nil || err2 != nil || page < 1 || limit < 1 {
			c.JSON(400, gin.H{"error": "invalid page or limit"})
			return
		}

		offset := (page - 1) * limit

		var res []ForChoose

		if err := dbInstant.Table("foods").
			Select("foods.name as food_name").
			Joins("JOIN food_ingredients ON food_ingredients.food_id = foods.id").
			Joins("JOIN ingredients ON ingredients.id = food_ingredients.ingredient_id").
			Group("foods.id").
			Having("SUM(CASE WHEN ingredients.name = ? THEN 0 ELSE 1 END) = 0", "miso").
			Limit(limit).
			Offset(offset).
			Scan(&res).Error; err != nil {

			c.JSON(500, gin.H{"error": "Failed to fetch food count"})
			return

		}
		c.JSON(200, gin.H{
			"page":  page,
			"limit": limit,
			"data":  res,
		})
	}
}

func UpdateFood(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req UpdatesFoodRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"message": "cannot input", "eror": err.Error()})
			return
		}

		var food_update model.Food = model.Food{
			Name:       req.Name,
			Price:      req.Price,
			CategoryID: req.CategoryID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := dbInstant.Where("id = ?", dataID).Updates(&food_update).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to get foods", "error": err.Error()})
			return
		}

		if err := dbInstant.Model(&model.Food_Ingredient{}).Where("food_id = ?", dataID).Delete(&model.Food_Ingredient{}).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to clear old ingredients", "error": err.Error()})
			return
		}
		var fi_update []model.Food_Ingredient = []model.Food_Ingredient{}
		var foodId, err1 = strconv.Atoi(dataID)
		if err1 != nil {
			c.JSON(400, gin.H{"message": "Failed to extract food id", "error": err1.Error()})
		}
		for _, ingredientID := range req.IngredientsID {
			fi_update = append(fi_update, model.Food_Ingredient{
				FoodID:       foodId,
				IngredientID: ingredientID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			})
		}
		food_ingredient_create_res := dbInstant.Model(&model.Food_Ingredient{}).Create(&fi_update)
		if food_ingredient_create_res.Error != nil {
			c.JSON(400, gin.H{"message": "Failed to create food ingredients", "error": food_ingredient_create_res.Error.Error()})
			return
		}
		c.JSON(200, gin.H{})
	}
}

func DeleteFood(dbInstant *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		dataID := c.Param("id")

		var req model.Food

		if err := dbInstant.Where("id = ?", dataID).Delete(&req).Error; err != nil {
			c.JSON(400, gin.H{"message": "Failed to get foods", "error": err.Error()})
			return
		}

		c.JSON(200, dataID)
	}
}
