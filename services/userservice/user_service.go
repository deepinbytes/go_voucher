package userservice

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/user"

	"github.com/deepinbytes/go_voucher/repositories/userrepo"
)

// UserService interface
type UserService interface {
	GetByID(id uint) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	ListAll() ([]*user.User, error)
	Create(*user.User) error
	Update(*user.User) error
}

type userService struct {
	Repo userrepo.Repo
}

func (us *userService) ListAll() ([]*user.User, error) {
	users, err := us.Repo.ListAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// NewUserService will instantiate User Service
func NewUserService(
	repo userrepo.Repo,
) UserService {

	return &userService{
		Repo: repo,
	}
}

func (us *userService) GetByID(id uint) (*user.User, error) {
	if id == 0 {
		return nil, errors.New("id param is required")
	}
	user, err := us.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) GetByEmail(email string) (*user.User, error) {
	if email == "" {
		return nil, errors.New("email(string) is required")
	}
	user, err := us.Repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) Create(user *user.User) error {
	return us.Repo.Create(user)
}

func (us *userService) Update(user *user.User) error {
	return us.Repo.Update(user)
}
