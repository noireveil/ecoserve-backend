package usecase

import (
	"errors"

	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
)

type CompleteOrderRequest struct {
	PhotoURL     string  `json:"photo_url"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	DeviceWeight float64 `json:"device_weight"`
	Category     string  `json:"category"`
	DistanceKm   float64 `json:"distance_km"`
	ServiceFee   float64 `json:"service_fee"`
}

type OrderUsecase interface {
	CreateOrder(order *domain.Order) error
	GetUserOrders(userID string) ([]domain.Order, error)
	GetOrderByID(orderID string) (*domain.Order, error)
	CompleteOrder(orderID string, req CompleteOrderRequest) error
	GetIncomingOrders() ([]domain.Order, error)
	AcceptOrder(orderID string, userID string) error
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

func (u *orderUsecase) GetUserOrders(userID string) ([]domain.Order, error) {
	return u.orderRepo.FindByUserID(userID)
}

func (u *orderUsecase) GetOrderByID(orderID string) (*domain.Order, error) {
	return u.orderRepo.FindByID(orderID)
}

func (u *orderUsecase) GetIncomingOrders() ([]domain.Order, error) {
	return u.orderRepo.FindIncomingOrders()
}

func (u *orderUsecase) AcceptOrder(orderID string, userID string) error {
	return u.orderRepo.AcceptOrder(orderID, userID)
}

func (u *orderUsecase) CompleteOrder(orderID string, req CompleteOrderRequest) error {
	if req.PhotoURL == "" {
		return errors.New("lapisan anti-fraud: bukti foto wajib dilampirkan")
	}
	if req.Longitude == 0 || req.Latitude == 0 {
		return errors.New("lapisan anti-fraud: verifikasi geospasial wajib disertakan")
	}
	if req.ServiceFee <= 0 {
		return errors.New("validasi gagal: teknisi harus memasukkan nominal biaya jasa riil yang disepakati")
	}

	totalFee := req.ServiceFee
	platformFee := totalFee * 0.10 // Take Rate 10% untuk EcoServe
	netFee := totalFee - platformFee

	eWasteSaved := u.calculateImpactMetrics(req.DeviceWeight, req.Category, req.DistanceKm)

	return u.orderRepo.CompleteWithAntiFraud(orderID, req.PhotoURL, req.Longitude, req.Latitude, eWasteSaved, totalFee, platformFee, netFee)
}

const EFTransportMotorcycle = 0.103

func (u *orderUsecase) getEFVirgin(category string) float64 {
	switch category {
	case "Smartphone", "IT & Gadget":
		return 70.0
	case "Laptop/PC":
		return 50.0
	case "Home Appliances", "Pendingin & Komersial":
		return 30.0
	default:
		return 30.0
	}
}

func (u *orderUsecase) calculateImpactMetrics(deviceWeightKg float64, category string, distanceKm float64) float64 {
	efVirgin := u.getEFVirgin(category)

	emisiProduksi := deviceWeightKg * efVirgin
	emisiTransportasi := distanceKm * EFTransportMotorcycle

	peTotal := emisiProduksi - emisiTransportasi

	if peTotal > 0 {
		return peTotal
	}
	return 0.0
}
