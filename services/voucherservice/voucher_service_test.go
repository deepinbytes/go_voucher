package voucherservice

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	t.Run("Get a voucher", func(t *testing.T) {
		expected := &voucher.Voucher{
			Code: "Test",
		}

		voucherRepo := new(repoMock)
		u := NewVoucherService(voucherRepo)
		voucherRepo.On("GetByID", testID100).Return(expected, nil)

		result, _ := u.GetByID(testID100)

		assert.EqualValues(t, expected, result)
	})

	t.Run("Get error if id is 0", func(t *testing.T) {
		expected := errors.New("id param is required")

		voucherRepo := new(repoMock)
		u := NewVoucherService(voucherRepo)

		result, err := u.GetByID(0)

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})

	t.Run("Get error if it has error", func(t *testing.T) {
		expected := errors.New("Nop")

		voucherRepo := new(repoMock)
		u := NewVoucherService(voucherRepo)
		voucherRepo.On("GetByID", testID10).Return(&voucher.Voucher{}, expected)

		result, err := u.GetByID(testID10)

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})
}

func TestRedeemCode(t *testing.T) {
	t.Run("Redeem a voucher", func(t *testing.T) {
		expected := &voucher.Voucher{
			Code: "Test",
		}

		voucherRepo := new(repoMock)

		u := NewVoucherService(voucherRepo)
		voucherRepo.On("UseCode", testName).Return(expected, nil)

		result, _ := u.UseCode(testName)

		assert.EqualValues(t, expected, result)
	})

	t.Run("Get error if offer is empty", func(t *testing.T) {
		expected := errors.New("Code(string) is required")

		voucherRepo := new(repoMock)

		u := NewVoucherService(voucherRepo)

		result, err := u.UseCode("")

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})

	t.Run("Get error if it has error", func(t *testing.T) {
		expected := errors.New("Nop")

		voucherRepo := new(repoMock)

		u := NewVoucherService(voucherRepo)
		voucherRepo.On("UseCode", testName).Return(&voucher.Voucher{}, expected)

		result, err := u.UseCode(testName)

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Create a voucher", func(t *testing.T) {
		offer := &voucher.Voucher{
			Code: "test",
		}

		voucherRepo := new(repoMock)

		u := NewVoucherService(voucherRepo)
		voucherRepo.On("Create", offer).Return(nil)

		result := u.Create(offer)

		assert.Nil(t, result)
	})

	t.Run("Create a voucher fails", func(t *testing.T) {
		err := errors.New(("oops"))
		offer := &voucher.Voucher{
			Code: "test",
		}

		voucherRepo := new(repoMock)

		u := NewVoucherService(voucherRepo)

		voucherRepo.On("Create", offer).Return(err)
		result := u.Create(offer)

		assert.EqualValues(t, result, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update a voucher", func(t *testing.T) {
		usr := &voucher.Voucher{
			Code: "Test",
		}

		voucherRepo := new(repoMock)

		u := NewVoucherService(voucherRepo)
		voucherRepo.On("Update", usr).Return(nil)

		result := u.Update(usr)

		assert.Nil(t, result)
	})

	t.Run("Update a offer fails", func(t *testing.T) {
		err := errors.New(("oops"))
		usr := &voucher.Voucher{
			Code: "test",
		}

		voucherRepo := new(repoMock)

		u := NewVoucherService(voucherRepo)
		voucherRepo.On("Update", usr).Return(err)

		result := u.Update(usr)

		assert.EqualValues(t, result, err)
	})
}
