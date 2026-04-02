package usecase

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
)

type OrderUsecase interface {
	CreateOrder(order *domain.Order) error
	CompleteOrder(orderID string, deviceCategory string) error
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{orderRepo}
}

func (u *orderUsecase) CreateOrder(order *domain.Order) error {
	return u.orderRepo.Create(order)
}

func (u *orderUsecase) CompleteOrder(orderID string, deviceCategory string) error {
	eWasteSaved := u.calculateImpactMetrics(deviceCategory)
	return u.orderRepo.UpdateStatus(orderID, "COMPLETED", eWasteSaved)
}

func (u *orderUsecase) calculateImpactMetrics(category string) float64 {
	metricsMap := map[string]float64{
		"Pendingin & Komersial": 45.50,
		"Home Appliances":       22.00,
		"IT & Gadget":           1.25,
	}

	if weight, exists := metricsMap[category]; exists {
		return weight
	}

	return 2.50
}
