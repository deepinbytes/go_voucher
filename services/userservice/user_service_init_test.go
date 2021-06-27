package userservice

import (
	"github.com/deepinbytes/go_voucher/domain/user"
	"github.com/stretchr/testify/mock"
)

var (
	testID10  = uint(10)
	testID100 = uint(100)
	testEmail = "test@cc.cc"
)

type repoMock struct {
	mock.Mock
}

func (repo *repoMock) ListAll() ([]*user.User, error) {
	args := repo.Called()
	return args.Get(0).([]*user.User), args.Error(1)
}

func (repo *repoMock) GetByID(id uint) (*user.User, error) {
	args := repo.Called(id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (repo *repoMock) GetByEmail(email string) (*user.User, error) {
	args := repo.Called(email)
	return args.Get(0).(*user.User), args.Error(1)
}

func (repo *repoMock) Create(user *user.User) error {
	args := repo.Called(user)
	return args.Error(0)
}

func (repo *repoMock) Update(user *user.User) error {
	args := repo.Called(user)
	return args.Error(0)
}
