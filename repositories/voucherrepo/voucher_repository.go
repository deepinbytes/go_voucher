package voucherrepo

import (
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"github.com/jinzhu/gorm"
)

// Repo interface
type Repo interface {
	GetByID(id uint) (*voucher.Voucher, error)
	UseCode(name string) (*voucher.Voucher, error)
	Create(voucher *voucher.Voucher) error
	Update(voucher *voucher.Voucher) error
}

type voucherRepo struct {
	db *gorm.DB
}

// NewUserRepo will instantiate User Repository
func NewVoucherRepo(db *gorm.DB) Repo {
	return &voucherRepo{
		db: db,
	}
}

func (u *voucherRepo) GetByID(id uint) (*voucher.Voucher, error) {
	var voucher voucher.Voucher
	if err := u.db.First(&voucher, id).Error; err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (u *voucherRepo) UseCode(name string) (*voucher.Voucher, error) {
	var v voucher.Voucher
	if err := u.db.First(&v, "vouchers.code = ?", name).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (u *voucherRepo) Create(voucher *voucher.Voucher) error {
	return u.db.Create(voucher).Error
}

func (u *voucherRepo) Update(voucher *voucher.Voucher) error {
	return u.db.Save(voucher).Error
}
