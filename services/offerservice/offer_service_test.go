package offerservice

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/offer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	t.Run("Get a user", func(t *testing.T) {
		expected := &offer.Offer{
			Name: "Test",
		}

		offerRepo := new(repoMock)
		u := NewOfferService(offerRepo)
		offerRepo.On("GetByID", testID100).Return(expected, nil)

		result, _ := u.GetByID(testID100)

		assert.EqualValues(t, expected, result)
	})

	t.Run("Get error if id is 0", func(t *testing.T) {
		expected := errors.New("id param is required")

		offerRepo := new(repoMock)
		u := NewOfferService(offerRepo)

		result, err := u.GetByID(0)

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})

	t.Run("Get error if it has error", func(t *testing.T) {
		expected := errors.New("Nop")

		offerRepo := new(repoMock)
		u := NewOfferService(offerRepo)
		offerRepo.On("GetByID", testID10).Return(&offer.Offer{}, expected)

		result, err := u.GetByID(testID10)

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})
}

func TestGetByEmail(t *testing.T) {
	t.Run("Get a offer", func(t *testing.T) {
		expected := &offer.Offer{
			Name: "Test",
		}

		offerRepo := new(repoMock)

		u := NewOfferService(offerRepo)
		offerRepo.On("GetByName", testName).Return(expected, nil)

		result, _ := u.GetByName(testName)

		assert.EqualValues(t, expected, result)
	})

	t.Run("Get error if offer is empty", func(t *testing.T) {
		expected := errors.New("Name(string) is required")

		offerRepo := new(repoMock)

		u := NewOfferService(offerRepo)

		result, err := u.GetByName("")

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})

	t.Run("Get error if it has error", func(t *testing.T) {
		expected := errors.New("Nop")

		offerRepo := new(repoMock)

		u := NewOfferService(offerRepo)
		offerRepo.On("GetByName", testName).Return(&offer.Offer{}, expected)

		result, err := u.GetByName(testName)

		assert.Nil(t, result)
		assert.EqualValues(t, expected, err)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Create a offer", func(t *testing.T) {
		offer := &offer.Offer{
			Name: "test",
		}

		offerRepo := new(repoMock)

		u := NewOfferService(offerRepo)
		offerRepo.On("Create", offer).Return(nil)

		result := u.Create(offer)

		assert.Nil(t, result)
	})

	t.Run("Create a offer fails", func(t *testing.T) {
		err := errors.New(("oops"))
		offer := &offer.Offer{
			Name: "test",
		}

		offerRepo := new(repoMock)

		u := NewOfferService(offerRepo)

		offerRepo.On("Create", offer).Return(err)
		result := u.Create(offer)

		assert.EqualValues(t, result, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update a offer", func(t *testing.T) {
		usr := &offer.Offer{
			Name: "Test",
		}

		offerRepo := new(repoMock)

		u := NewOfferService(offerRepo)
		offerRepo.On("Update", usr).Return(nil)

		result := u.Update(usr)

		assert.Nil(t, result)
	})

	t.Run("Update a offer fails", func(t *testing.T) {
		err := errors.New(("oops"))
		usr := &offer.Offer{
			Name: "test",
		}

		offerRepo := new(repoMock)

		u := NewOfferService(offerRepo)
		offerRepo.On("Update", usr).Return(err)

		result := u.Update(usr)

		assert.EqualValues(t, result, err)
	})
}
