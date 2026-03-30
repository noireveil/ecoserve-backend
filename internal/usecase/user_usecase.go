package usecase

import (
	"errors"

	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
	"gorm.io/gorm"
)

type UserUsecase interface {
	LoginOrRegister(fullName, whatsapp string) (*domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) LoginOrRegister(fullName, whatsapp string) (*domain.User, error) {
	user, err := u.userRepo.FindByWhatsApp(whatsapp)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := &domain.User{
				FullName:       fullName,
				WhatsAppNumber: whatsapp,
				Role:           "customer",
			}
			if errCreate := u.userRepo.Create(newUser); errCreate != nil {
				return nil, errCreate
			}
			return newUser, nil
		}
		return nil, err
	}

	return user, nil
}
