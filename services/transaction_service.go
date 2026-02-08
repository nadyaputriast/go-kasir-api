package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items, useLock)
}

func (s *TransactionService) GetReport(start, end string) (*models.SalesSummary, error) {
	var startDate, endDate time.Time
	var err error

	if start == "" || end == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	} else {
		layout := "2006-01-02"
		startDate, err = time.Parse(layout, start)
		if err != nil {
			return nil, err
		}

		endDate, err = time.Parse(layout, end)
		if err != nil {
			return nil, err
		}
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	return s.repo.GetSalesSummary(startDate, endDate)
}
