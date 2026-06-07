package model

import (
	"time"

	"gorm.io/gorm"
)

type Food_Ingredient struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement:true"`
	FoodID       int            `json:"food_id"`
	IngredientID int            `json:"ingredient_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}
