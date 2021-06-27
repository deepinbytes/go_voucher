package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// NOTE: Mocked services are in './offer_controller_setup_test.go'

// Output of HTTP Response Body structure
type outputVoucher struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data voucher.Voucher `json:"data"`
}

func TestVoucherController(t *testing.T) {

	// Setup router + offer controller
	os := &offerSvc{}
	us := &userSvc{}
	vs := &voucherSvc{}
	voucherCtl := NewVoucherController(vs, us, os)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/voucher/:id", voucherCtl.GetByID)

	// Using router version
	t.Run("GetByID", func(t *testing.T) {
		t.Run("Get a voucher", func(t *testing.T) {
			// Make HTTP Request to the testing endpoint
			w := performRequest(router, "GET", "/voucher/1")

			// Check statusCode
			assert.Equal(t, http.StatusOK, w.Code)

			// JSON to struct
			resBody := outputVoucher{}
			json.NewDecoder(w.Body).Decode(&resBody)

			// Expected HTTP Response body structure
			expectedResBody := Response{
				Code: http.StatusOK,
				Msg:  "ok",
				Data: *voucher1,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})

		t.Run("Fails to get a voucher without valid id", func(t *testing.T) {
			w := performRequest(router, "GET", "/voucher/b")

			assert.Equal(t, http.StatusBadRequest, w.Code)

			resBody := failedOutput{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusBadRequest,
				Msg:  "voucher id should be a number",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})

		t.Run("Fails to get a voucher (not found))", func(t *testing.T) {
			w := performRequest(router, "GET", "/voucher/10")

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

		t.Run("Fails to get a voucher (something went wrong))", func(t *testing.T) {
			w := performRequest(router, "GET", "/voucher/100")

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
				"code":     "test",
				"user_id":  1,
				"offer_id": 1,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Mock Request body
			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/voucher/Create", bytes.NewBuffer(payload))
			// request.Header.Set("content-type", "application/json")
			// router.ServeHTTP(w, request)
			c.Request = request

			voucherCtl.Create(c)

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
			request := httptest.NewRequest("POST", "/voucher/create", nil)
			c.Request = request

			voucherCtl.Create(c)

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

		t.Run("Fails to create voucher", func(t *testing.T) {
			reqBody := map[string]interface{}{
				"code":     "existing_code",
				"user_id":  45,
				"offer_id": 12,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/voucher/create", bytes.NewBuffer(payload))
			c.Request = request

			voucherCtl.Create(c)

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

	t.Run("Redeem", func(t *testing.T) {
		t.Run("Offer not available", func(t *testing.T) {
			reqBody := map[string]interface{}{
				"code": "TEST1",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/voucher/redeem", bytes.NewBuffer(payload))
			c.Request = request

			voucherCtl.Redeem(c)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)
			expectedResBody := Response{
				Code: 500,
				Msg:  "Nop",
				Data: "Offer Not Available Anymore",
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})
	})

}
