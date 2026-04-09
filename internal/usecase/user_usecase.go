package usecase

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
	"github.com/noireveil/ecoserve-backend/pkg/utils"
	"gorm.io/gorm"
)

type UserUsecase interface {
	RequestOTP(fullName, email string) error
	VerifyOTP(email, code string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) RequestOTP(fullName, email string) error {
	_, err := u.userRepo.FindByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := &domain.User{
				FullName: fullName,
				Email:    email,
				Role:     "customer",
			}
			if errCreate := u.userRepo.Create(newUser); errCreate != nil {
				return errCreate
			}
		} else {
			return err
		}
	}

	otpCode, _ := utils.GenerateOTP()
	expiresAt := time.Now().Add(time.Minute * 5)

	if errUpdate := u.userRepo.UpdateOTP(email, otpCode, expiresAt); errUpdate != nil {
		return errUpdate
	}

	go func(target, code string) {
		if errSend := utils.SendEmailOTP(target, code); errSend != nil {
			log.Printf("Gagal mengirim OTP ke %s: %v\n", target, errSend)
		} else {
			log.Printf("OTP berhasil dikirim ke email: %s\n", target)
		}
	}(email, otpCode)

	return nil
}

func (u *userUsecase) VerifyOTP(email, code string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	cleanCode := strings.TrimSpace(code)

	if user.OTPCode != cleanCode {
		return nil, errors.New("kode OTP tidak valid")
	}

	if time.Now().After(user.OTPExpiresAt) {
		return nil, errors.New("kode OTP telah kadaluarsa")
	}

	_ = u.userRepo.UpdateOTP(email, "", time.Now())

	return user, nil
}

func (u *userUsecase) GetUserByID(id string) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}
