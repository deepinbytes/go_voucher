package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deepinbytes/go_voucher/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// NOTE: Mocked services are in './user_controller_setup_test.go'

// Output of HTTP Response Body structure
type output struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data user.User `json:"data"`
}

type failedOutput struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestUserController(t *testing.T) {

	// Setup router + user controller
	us := &userSvc{}
	userCtl := NewUserController(us)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", userCtl.GetByID)

	// Using router version
	t.Run("GetByID", func(t *testing.T) {
		t.Run("Get a user", func(t *testing.T) {
			// Make HTTP Request to the testing endpoint
			w := performRequest(router, "GET", "/users/1")

			// Check statusCode
			assert.Equal(t, http.StatusOK, w.Code)

			// JSON to struct
			resBody := output{}
			json.NewDecoder(w.Body).Decode(&resBody)

			// Expected HTTP Response body structure
			expectedResBody := Response{
				Code: http.StatusOK,
				Msg:  "ok",
				Data: *alice,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})

		t.Run("Fails to get a user without valid id", func(t *testing.T) {
			w := performRequest(router, "GET", "/users/b")

			assert.Equal(t, http.StatusBadRequest, w.Code)

			resBody := failedOutput{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusBadRequest,
				Msg:  "user id should be a number",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody.Code, resBody.Code)
			assert.EqualValues(t, expectedResBody.Msg, resBody.Msg)
			assert.EqualValues(t, expectedResBody.Data, resBody.Data)
		})

		t.Run("Fails to get a user (not found))", func(t *testing.T) {
			w := performRequest(router, "GET", "/users/10")

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

		t.Run("Fails to get a user (something went wrong))", func(t *testing.T) {
			w := performRequest(router, "GET", "/users/100")

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

	t.Run("Register", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			reqBody := map[string]string{
				"email":    "alice@cc.cc",
				"password": "123test",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Mock Request body
			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/register", bytes.NewBuffer(payload))
			// request.Header.Set("content-type", "application/json")
			// router.ServeHTTP(w, request)
			c.Request = request

			userCtl.Register(c)

			assert.Equal(t, http.StatusOK, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusOK,
				Msg:  "ok",
				Data: map[string]interface{}{
					"id":        float64(0),
					"email":     "alice@cc.cc",
					"firstName": "",
					"lastName":  "",
					"vouchers":  interface{}(nil),
				},
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

		t.Run("Invalid payload", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Mock Request body
			request := httptest.NewRequest("POST", "/register", nil)
			c.Request = request

			userCtl.Register(c)

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

		t.Run("Fails to create user", func(t *testing.T) {
			reqBody := map[string]string{
				"email":    "bob@cc.cc",
				"password": "123test",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("POST", "/register", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Register(c)

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
			reqBody := map[string]string{
				"email":     "alice@cc.cc",
				"firstName": "alice",
				"lastName":  "smith",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user_id", uint(1))

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("PUT", "/profile", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Update(c)

			assert.Equal(t, http.StatusOK, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusOK,
				Msg:  "ok",
				Data: map[string]interface{}{
					"id":        float64(1),
					"email":     "alice@cc.cc",
					"firstName": "alice",
					"lastName":  "smith",
					"vouchers":  interface{}(nil),
				},
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

		t.Run("Fails to get user from db", func(t *testing.T) {
			reqBody := map[string]string{
				"email":     "alice@cc.cc",
				"firstName": "alice",
				"lastName":  "smith",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user_id", uint(0))

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("PUT", "/profile", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Update(c)

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

		t.Run("Invalid User ID", func(t *testing.T) {
			reqBody := map[string]string{
				"email":     "alice@cc.cc",
				"firstName": "alice",
				"lastName":  "smith",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user_id", uint(3))

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("PUT", "/profile", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Update(c)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody, resBody)
		})

		t.Run("Fails to update", func(t *testing.T) {
			reqBody := map[string]string{
				"email":     "bob@cc.cc",
				"firstName": "alice",
				"lastName":  "smith",
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("user_id", uint(1))

			payload, _ := json.Marshal(reqBody)
			request := httptest.NewRequest("PUT", "/profile", bytes.NewBuffer(payload))
			c.Request = request

			userCtl.Update(c)

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			resBody := Response{}
			json.NewDecoder(w.Body).Decode(&resBody)

			expectedResBody := Response{
				Code: http.StatusInternalServerError,
				Msg:  "Nop",
				Data: nil,
			}

			assert.EqualValues(t, expectedResBody, resBody)

			// Reset
			alice.Email = "alice@cc.cc"
		})
	})

}
