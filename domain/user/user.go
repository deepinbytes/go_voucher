package user

import (
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"github.com/jinzhu/gorm"
)

// User domain model
type User struct {
	gorm.Model
	FirstName string            `gorm:"size:255" `
	LastName  string            `gorm:"size:255"`
	Email     string            `gorm:"NOT NULL; UNIQUE_INDEX"`
	Voucher   []voucher.Voucher `gorm:"foreignKey:UserID"`
}
