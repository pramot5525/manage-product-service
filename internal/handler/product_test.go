package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"prm-product/internal/model"
	servicemocks "prm-product/internal/service/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		setupMock      func(s *servicemocks.ProductService)
		expectStatus   int
		expectSuccess  bool
		expectError    string
		expectDataID   int64
		expectDataName string
	}{
		{
			name:         "invalid json body",
			body:         `{"name":`,
			expectStatus: http.StatusBadRequest,
			expectError:  "invalid_request_body",
		},
		{
			name:         "validation error",
			body:         `{"price":999.99}`,
			expectStatus: http.StatusBadRequest,
			expectError:  "validation_error",
		},
		{
			name: "service error",
			body: `{"name":"iPhone 15","price":999.99}`,
			setupMock: func(s *servicemocks.ProductService) {
				s.EXPECT().Create(mock.AnythingOfType("*model.ProductRequest")).Return(int64(0), errors.New("create failed"))
			},
			expectStatus: http.StatusInternalServerError,
			expectError:  "failed_to_create_product",
		},
		{
			name: "success",
			body: `{"name":"iPhone 15","price":999.99}`,
			setupMock: func(s *servicemocks.ProductService) {
				s.EXPECT().Create(mock.MatchedBy(func(req *model.ProductRequest) bool {
					return req != nil && req.Name == "iPhone 15" && req.Price == 999.99
				})).Return(int64(101), nil)
			},
			expectStatus:   http.StatusCreated,
			expectSuccess:  true,
			expectDataID:   101,
			expectDataName: "iPhone 15",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			serviceMock := servicemocks.NewProductService(t)
			if tc.setupMock != nil {
				tc.setupMock(serviceMock)
			}

			handler := NewProductHandler(serviceMock)
			app := fiber.New()
			app.Post("/products/create", handler.CreateProduct)

			req := httptest.NewRequest(http.MethodPost, "/products/create", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectStatus, resp.StatusCode)

			var payload model.APIResponse
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
			require.Equal(t, tc.expectSuccess, payload.Successful)

			if tc.expectError != "" {
				require.Equal(t, tc.expectError, payload.ErrorCode)
			}

			if tc.expectDataID != 0 {
				require.NotNil(t, payload.Data)
				require.Equal(t, tc.expectDataID, payload.Data.ID)
				require.Equal(t, tc.expectDataName, payload.Data.Name)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	testCases := []struct {
		name          string
		pathID        string
		body          string
		setupMock     func(s *servicemocks.ProductService)
		expectStatus  int
		expectSuccess bool
		expectError   string
	}{
		{
			name:         "invalid json body",
			pathID:       "1",
			body:         `{"price":`,
			expectStatus: http.StatusBadRequest,
			expectError:  "invalid_request_body",
		},
		{
			name:         "invalid id",
			pathID:       "abc",
			body:         `{"price":1000}`,
			expectStatus: http.StatusBadRequest,
			expectError:  "invalid_id",
		},
		{
			name:         "validation error from unknown field",
			pathID:       "1",
			body:         `{"foo":"bar"}`,
			expectStatus: http.StatusBadRequest,
			expectError:  "validation_error",
		},
		{
			name:   "service error",
			pathID: "15",
			body:   `{"price":1200}`,
			setupMock: func(s *servicemocks.ProductService) {
				s.EXPECT().Update(int64(15), mock.MatchedBy(func(updateData map[string]interface{}) bool {
					price, ok := updateData["price"].(float64)
					return ok && price == 1200
				})).Return(errors.New("update failed"))
			},
			expectStatus: http.StatusInternalServerError,
			expectError:  "failed_to_update_product",
		},
		{
			name:   "success with nullable sale_price",
			pathID: "20",
			body:   `{"sale_price":null}`,
			setupMock: func(s *servicemocks.ProductService) {
				s.EXPECT().Update(int64(20), mock.MatchedBy(func(updateData map[string]interface{}) bool {
					value, exists := updateData["sale_price"]
					return exists && value == nil
				})).Return(nil)
			},
			expectStatus:  http.StatusOK,
			expectSuccess: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			serviceMock := servicemocks.NewProductService(t)
			if tc.setupMock != nil {
				tc.setupMock(serviceMock)
			}

			handler := NewProductHandler(serviceMock)
			app := fiber.New()
			app.Patch("/products/update/:id", handler.UpdateProduct)

			req := httptest.NewRequest(http.MethodPatch, "/products/update/"+tc.pathID, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectStatus, resp.StatusCode)

			var payload model.APIResponse
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
			require.Equal(t, tc.expectSuccess, payload.Successful)

			if tc.expectError != "" {
				require.Equal(t, tc.expectError, payload.ErrorCode)
			}
		})
	}
}