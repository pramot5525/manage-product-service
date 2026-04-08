package service

import (
	"errors"
	"testing"

	"manage-product-service/internal/model"
	repositorymocks "manage-product-service/internal/repository/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProductServiceCreate(t *testing.T) {
	testCases := []struct {
		name       string
		req        *model.ProductRequest
		mockReturn int64
		mockErr    error
		expectID   int64
		expectErr  error
		assertArg  func(*model.ProductRequest) func(*model.Product) bool
	}{
		{
			name: "success",
			req: func() *model.ProductRequest {
				description := "Smartphone"
				salePrice := 899.99
				return &model.ProductRequest{
					Name:        "iPhone 15",
					Description: &description,
					SalePrice:   &salePrice,
					Price:       999.99,
				}
			}(),
			mockReturn: 101,
			expectID:   101,
			assertArg: func(req *model.ProductRequest) func(*model.Product) bool {
				return func(product *model.Product) bool {
					return product != nil &&
						product.ID == 0 &&
						product.Name == req.Name &&
						product.Description != nil &&
						*product.Description == *req.Description &&
						product.SalePrice != nil &&
						*product.SalePrice == *req.SalePrice &&
						product.Price == req.Price
				}
			},
		},
		{
			name: "repository error",
			req: &model.ProductRequest{
				Name:  "iPhone 15",
				Price: 999.99,
			},
			mockErr:   errors.New("insert failed"),
			expectErr: errors.New("insert failed"),
			assertArg: func(_ *model.ProductRequest) func(*model.Product) bool {
				return func(product *model.Product) bool { return product != nil }
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := repositorymocks.NewProductRepository(t)
			service := NewProductService(repo)

			repo.EXPECT().Create(mock.MatchedBy(tc.assertArg(tc.req))).Return(tc.mockReturn, tc.mockErr)

			id, err := service.Create(tc.req)

			if tc.expectErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectID, id)
		})
	}
}

func TestProductServiceUpdate(t *testing.T) {
	testCases := []struct {
		name       string
		id         int64
		updateData map[string]interface{}
		mockErr    error
		expectErr  error
	}{
		{
			name: "success",
			id:   15,
			updateData: map[string]interface{}{
				"name":       "iPhone 15 Pro",
				"sale_price": nil,
			},
		},
		{
			name: "repository error",
			id:   20,
			updateData: map[string]interface{}{
				"price": 1200.0,
			},
			mockErr:   errors.New("update failed"),
			expectErr: errors.New("update failed"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := repositorymocks.NewProductRepository(t)
			service := NewProductService(repo)

			repo.EXPECT().Update(tc.id, tc.updateData).Return(tc.mockErr)

			err := service.Update(tc.id, tc.updateData)

			if tc.expectErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
