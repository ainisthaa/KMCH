package model

import (
	"time"

	"gorm.io/gorm"
)

type Ingredient struct {
	ID              int               `json:"id" gorm:"primaryKey"`
	Name            string            `json:"name"`
	Food_Ingredient []Food_Ingredient `gorm:"foreignKey:IngredientID"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `json:"deleted_at"`
}
