package repository

import (
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	CreateAndUpdateRating(review *domain.Review) error
	CheckExistsByOrderID(orderID string) (bool, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db}
}

func (r *reviewRepository) CheckExistsByOrderID(orderID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Review{}).Where("order_id = ?", orderID).Count(&count).Error
	return count > 0, err
}

func (r *reviewRepository) CreateAndUpdateRating(review *domain.Review) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(review).Error; err != nil {
			return err
		}

		var result struct {
			AvgRating float64
		}
		if err := tx.Model(&domain.Review{}).
			Select("COALESCE(AVG(rating), 0) as avg_rating").
			Where("technician_id = ?", review.TechnicianID).
			Scan(&result).Error; err != nil {
			return err
		}

		if err := tx.Model(&domain.Technician{}).
			Where("id = ?", review.TechnicianID).
			Update("rating", result.AvgRating).Error; err != nil {
			return err
		}

		return nil
	})
}
