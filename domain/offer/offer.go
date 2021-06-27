package offer

import (
	"github.com/jinzhu/gorm"
)

// User domain model
type Offer struct {
	gorm.Model
	Name               string `gorm:"NOT NULL; UNIQUE_INDEX" json:"name"`
	DiscountPercentage uint   `json:"discount_percentage"`
}
