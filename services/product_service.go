package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Create(product *models.Product) error {
	isDuplicate, err := s.repo.Exists(product.Name, product.Price, product.CategoryId)
	if err != nil {
		return err
	}
	if isDuplicate {
		return errors.New("a product with the same name, price, and category already exists")
	}

	if product.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if product.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	err = s.repo.Create(product)
	if err != nil {
		return err
	}

	fullData, err := s.repo.GetByID(product.ID)
	if err == nil {
		*product = *fullData
	}

	return nil
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) error {
	existingProduct, err := s.repo.GetByID(product.ID)
	if err != nil {
		return err
	}

	if product.Name == existingProduct.Name &&
		product.Price == existingProduct.Price &&
		product.Stock == existingProduct.Stock &&
		product.CategoryId == existingProduct.CategoryId {
		return errors.New("no changes detected; the updated data is identical to the current data")
	}

	if product.Name == "" {
		product.Name = existingProduct.Name
	}
	if product.Price == 0 {
		product.Price = existingProduct.Price
	}
	if product.Stock == 0 {
		product.Stock = existingProduct.Stock
	}
	if product.CategoryId == 0 {
		product.CategoryId = existingProduct.CategoryId
	}

	err = s.repo.Update(product)
	if err != nil {
		return err
	}

	fullData, err := s.repo.GetByID(product.ID)
	if err == nil {
		*product = *fullData
	}

	return nil
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
