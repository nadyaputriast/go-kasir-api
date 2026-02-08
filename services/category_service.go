package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(name string) ([]models.Category, error) {
	return s.repo.GetAll(name)
}

func (s *CategoryService) Create(category *models.Category) error {
	isDuplicate, err := s.repo.Exists(category.Name, category.Description)
	if err != nil {
		return err
	}
	if isDuplicate {
		return errors.New("a category with the same name and description already exists")
	}

	return s.repo.Create(category)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(category *models.Category) error {
	existingCategory, err := s.repo.GetByID(category.ID)
	if err != nil {
		return err
	}

	if category.Name == existingCategory.Name &&
		category.Description == existingCategory.Description {
		return errors.New("no changes detected; the updated data is identical to the current data")
	}

	if category.Name == "" {
		category.Name = existingCategory.Name
	}
	if category.Description == "" {
		category.Description = existingCategory.Description
	}

	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
