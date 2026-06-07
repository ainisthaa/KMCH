package model

import (
	"time"

	"gorm.io/gorm"
)

type Food struct {
	ID              int               `json:"id" gorm:"primaryKey"`
	Name            string            `json:"name"`
	Price           uint              `json:"price"`
	Food_Ingredient []Food_Ingredient `gorm:"foreignKey:FoodID"`
	CategoryID      int               `json:"category_id"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"deleted_at"`
}
