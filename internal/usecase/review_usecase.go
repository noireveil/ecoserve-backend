package usecase

import (
	"errors"

	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
)

type SubmitReviewPayload struct {
	OrderID string
	UserID  string
	Rating  int
	Comment string
}

type ReviewUsecase interface {
	SubmitReview(payload SubmitReviewPayload) error
}

type reviewUsecase struct {
	reviewRepo repository.ReviewRepository
	orderRepo  repository.OrderRepository
}

func NewReviewUsecase(reviewRepo repository.ReviewRepository, orderRepo repository.OrderRepository) ReviewUsecase {
	return &reviewUsecase{reviewRepo, orderRepo}
}

func (u *reviewUsecase) SubmitReview(payload SubmitReviewPayload) error {
	if payload.Rating < 1 || payload.Rating > 5 {
		return errors.New("rating tidak valid, harus antara 1 hingga 5")
	}

	order, err := u.orderRepo.FindByID(payload.OrderID)
	if err != nil {
		return errors.New("pesanan tidak ditemukan")
	}

	if order.CustomerID.String() != payload.UserID {
		return errors.New("akses ditolak: anda bukan pemilik pesanan ini")
	}

	if order.Status != domain.OrderStatusCompleted {
		return errors.New("ulasan hanya dapat diberikan pada pesanan yang telah selesai")
	}

	if order.TechnicianID == nil {
		return errors.New("integritas data gagal: pesanan ini tidak memiliki teknisi terdaftar")
	}

	exists, err := u.reviewRepo.CheckExistsByOrderID(payload.OrderID)
	if err != nil {
		return errors.New("gagal memvalidasi status ulasan")
	}
	if exists {
		return errors.New("anda sudah memberikan ulasan untuk pesanan ini")
	}

	review := &domain.Review{
		OrderID:      order.ID,
		CustomerID:   order.CustomerID,
		TechnicianID: *order.TechnicianID,
		Rating:       payload.Rating,
		Comment:      payload.Comment,
	}

	return u.reviewRepo.CreateAndUpdateRating(review)
}
