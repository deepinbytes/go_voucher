package voucher

import (
	"github.com/deepinbytes/go_voucher/domain/offer"
	"github.com/jinzhu/gorm"
	"time"
)

// User domain model
type Voucher struct {
	gorm.Model
	UsedAt     time.Time    `json:"used_at"`
	IsUsed     bool         `gorm:"default:false" json:"is_used"`
	Code       string       `gorm:"NOT NULL; UNIQUE_INDEX " json:"code"`
	OfferID    uint         `gorm:"foreignKey:OfferID" json:"offer_id"`
	UserID     uint         `json:"user_id "`
	ExpireTime time.Time    `json:"expiry_time"`
	Offer      *offer.Offer `gorm:"foreignKey:OfferID" json:"offer"`
}
