package userrepo

import (
	"github.com/deepinbytes/go_voucher/domain/user"

	"github.com/jinzhu/gorm"
)

// Repo interface
type Repo interface {
	GetByID(id uint) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Create(user *user.User) error
	Update(user *user.User) error
	ListAll() ([]*user.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func (u *userRepo) ListAll() ([]*user.User, error) {
	var users []*user.User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// NewUserRepo will instantiate User Repository
func NewUserRepo(db *gorm.DB) Repo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) GetByID(id uint) (*user.User, error) {
	var user user.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) GetByEmail(email string) (*user.User, error) {
	var user user.User
	if err := u.db.Preload("Voucher", "is_used = ?", "false").Preload("Voucher.Offer").
		Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) Create(user *user.User) error {
	return u.db.Create(user).Error
}

func (u *userRepo) Update(user *user.User) error {
	return u.db.Save(user).Error
}
