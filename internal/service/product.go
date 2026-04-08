package service

import (
	"manage-product-service/internal/model"
	"manage-product-service/internal/repository"
)

type ProductService interface {
	Create(*model.ProductRequest) (int64, error)
	Update(int64, map[string]interface{}) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) Create(product *model.ProductRequest) (int64, error) {
	formattedProduct := &model.Product{
		ID:          0,
		Name:        product.Name,
		Description: product.Description,
		SalePrice:   product.SalePrice,
		Price:       product.Price,
	}
	return s.repo.Create(formattedProduct)
}

func (s *productService) Update(id int64, updateData map[string]interface{}) error {
	if err := s.repo.Update(id, updateData); err != nil {
		return err
	}
	return nil
}
