package repository

import (
	"manage-product-service/internal/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(*model.Product) (int64, error)
	Update(int64, map[string]interface{}) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(product *model.Product) (int64, error) {
	if err := r.db.Create(product).Error; err != nil {
		return 0, err
	}

	return product.ID, nil
}

func (r *productRepository) Update(id int64, updateData map[string]interface{}) error {
	if len(updateData) == 0 {
		return nil
	}

	if err := r.db.Model(&model.Product{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}
