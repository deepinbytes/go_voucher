package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/deepinbytes/go_voucher/domain/offer"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// NOTE: Mocked services are in './offer_controller_setup_test.go'

// Output of HTTP Response Body structure
type outputOffer struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data offer.Offer `json:"data"`
}

func TestOfferController(t *testing.T) {

	// Setup router + offer controller
	os := &offerSvc{}
	us := &userSvc{}
	vs := &voucherSvc{}
	offerCtl := NewOfferController(os, us, vs)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/offer/:id", offerCtl.GetByID)

	// Using router version
	t.Run("GetByID", func(t *testing.T) {
		t.Run("Get a offer", func(t *testing.T) {
			// Make HTTP Request to the testing endpoint
			w := performRequest(router, "GET", "/offer/1")

			// Check statusCode
			assert.Equal(t, http.StatusOK, w.Code)

			// JSON to struct
			resBody := outputOffer{}
			json.NewDecoder(w.Body).Decode(&resBody)

			// Expected HTTP Response body structure
			expectedResBody := Response{
				Code: http.StatusOK,
				Msg:  "ok",
				Data: *of1,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})

		t.Run("Fails to get a offer without valid id", func(t *testing.T) {
			w := performRequest(router, "GET", "/offer/b")

			assert.Equal(t, http.StatusBadRequest, w.Code)

			resBody := failedOutput{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusBadRequest,
				Msg:  "offer id should be a number",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})

		t.Run("Fails to get a offer (not found))", func(t *testing.T) {
			w := performRequest(router, "GET", "/offer/10")

			assert.Equal(t, http.StatusNotFound, w.Code)

			resBody := failedOutput{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusNotFound,
				Msg:  "Record not found",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})

		t.Run("Fails to get a offer (something went wrong))", func(t *testing.T) {
			w := performRequest(router, "GET", "/offer/100")

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			resBody := failedOutput{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusInternalServerError,
				Msg:  "Ugh",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":                "test",
				"discount_percentage": 27,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Mock Request body
			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/offer/Create", bytes.NewBuffer(payload))
			// request.Header.Set("content-type", "application/json")
			// router.ServeHTTP(w, request)
			c.Request = request

			offerCtl.Create(c)

			assert.Equal(t, http.StatusOK, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: 0,
				Msg:  "",
				Data: interface{}(nil),
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

		t.Run("Invalid payload", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Mock Request body
			request := httptest.NewRequest("POST", "/offer/create", nil)
			c.Request = request

			offerCtl.Create(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusBadRequest,
				Msg:  "EOF",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

		t.Run("Fails to create offer", func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":                "new_offer",
				"discount_percentage": 45,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/offer/create", bytes.NewBuffer(payload))
			c.Request = request

			offerCtl.Create(c)

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusInternalServerError,
				Msg:  "Nop",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

	})

	t.Run("Update", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":                "offer1",
				"discount_percentage": 45,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("offer_id", uint(1))

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("PUT", "/offer/update", bytes.NewBuffer(payload))
			c.Request = request

			offerCtl.Update(c)

			assert.Equal(t, http.StatusOK, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusOK,
				Msg:  "ok",
				Data: map[string]interface{}{
					"id":                  float64(1),
					"name":                "offer1",
					"discount_percentage": float64(45),
				},
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

		t.Run("Fails to get offer from db", func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":                "offer1",
				"discount_percentage": 45,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("offer_id", uint(0))

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("PUT", "/offer/update", bytes.NewBuffer(payload))
			c.Request = request

			offerCtl.Update(c)

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusInternalServerError,
				Msg:  "Nop",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

		t.Run("Fails to update", func(t *testing.T) {
			reqBody := map[string]interface{}{
				"name":                "non_existing_offer",
				"discount_percentage": 90,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("offer_id", uint(1))

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("PUT", "/offer/update", bytes.NewBuffer(payload))
			c.Request = request

			offerCtl.Update(c)

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusInternalServerError,
				Msg:  "Nop",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody, resBody)

		})
	})

}
