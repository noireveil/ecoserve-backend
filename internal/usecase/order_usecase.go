package usecase

import (
	"errors"

	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
	"github.com/noireveil/ecoserve-backend/pkg/utils"
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
	GetUserOrders(userID string, limit, offset int) ([]domain.Order, error)
	GetOrderByID(orderID string) (*domain.Order, error)
	CompleteOrder(orderID string, req CompleteOrderRequest) error
	GetIncomingOrders(limit, offset int) ([]domain.Order, error)
	AcceptOrder(orderID string, userID string) error
	CancelOrder(orderID string, userID string) error
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
	techRepo  repository.TechnicianRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository, techRepo repository.TechnicianRepository) OrderUsecase {
	return &orderUsecase{orderRepo, techRepo}
}

func (u *orderUsecase) CreateOrder(order *domain.Order) error {
	if err := u.orderRepo.Create(order); err != nil {
		return err
	}

	go func(o domain.Order) {
		if o.TechnicianID != nil {
			tech, err := u.techRepo.FindByID(o.TechnicianID.String())
			if err == nil && tech.User.Email != "" {
				_ = utils.SendOrderNotificationEmail(tech.User.Email, tech.User.FullName, o.DeviceCategory, o.ProblemDescription)
			}
		} else {
			nearbyTechs, err := u.techRepo.FindNearby(o.CustomerLongitude, o.CustomerLatitude, 15)
			if err == nil {
				limit := 5
				if len(nearbyTechs) < limit {
					limit = len(nearbyTechs)
				}
				for i := 0; i < limit; i++ {
					if nearbyTechs[i].User.Email != "" {
						_ = utils.SendOrderNotificationEmail(nearbyTechs[i].User.Email, nearbyTechs[i].User.FullName, o.DeviceCategory, o.ProblemDescription)
					}
				}
			}
		}
	}(*order)

	return nil
}

func (u *orderUsecase) GetUserOrders(userID string, limit, offset int) ([]domain.Order, error) {
	return u.orderRepo.FindByUserID(userID, limit, offset)
}

func (u *orderUsecase) GetOrderByID(orderID string) (*domain.Order, error) {
	return u.orderRepo.FindByID(orderID)
}

func (u *orderUsecase) GetIncomingOrders(limit, offset int) ([]domain.Order, error) {
	return u.orderRepo.FindIncomingOrders(limit, offset)
}

func (u *orderUsecase) AcceptOrder(orderID string, userID string) error {
	return u.orderRepo.AcceptOrder(orderID, userID)
}

func (u *orderUsecase) CancelOrder(orderID string, userID string) error {
	return u.orderRepo.CancelOrder(orderID, userID)
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
	platformFee := totalFee * 0.10
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
