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
	var eWasteSaved float64

	switch deviceCategory {
	case "Pendingin & Komersial":
		eWasteSaved = 45.0
	case "Home Appliances":
		eWasteSaved = 25.0
	case "IT & Gadget":
		eWasteSaved = 2.5
	default:
		eWasteSaved = 5.0
	}

	return u.orderRepo.UpdateStatus(orderID, "COMPLETED", eWasteSaved)
}
