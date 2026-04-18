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

type ConsumerImpactDTO struct {
	TotalRepairs int     `json:"total_repairs" example:"3"`
	CO2AvoidedKg float64 `json:"co2_avoided_kg" example:"45.5"`
}

type UpdateProfileRequest struct {
	FullName          string  `json:"full_name" example:"EcoServe Tester"`
	ProfilePictureURL *string `json:"profile_picture_url" example:"https://storage.com/photo.jpg"`
}

type UserUsecase interface {
	RequestOTP(fullName, email string) error
	VerifyOTP(email, code string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	DeleteAccount(id string) error
	GetImpact(id string) (ConsumerImpactDTO, error)
	UpdateProfile(id string, req UpdateProfileRequest) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) RequestOTP(fullName, email string) error {
	user, err := u.userRepo.FindUnscopedByEmail(email)
	var targetName string

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			targetName = fullName

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
	} else {
		if user.DeletedAt.Valid {
			targetName = fullName
			_ = u.userRepo.HardDeleteTechnicianProfile(user.ID.String())
			if errRestore := u.userRepo.RestoreAndUpdate(email, fullName); errRestore != nil {
				return errRestore
			}
		} else {
			targetName = user.FullName
		}
	}

	otpCode, _ := utils.GenerateOTP()
	expiresAt := time.Now().Add(time.Minute * 5)

	if errUpdate := u.userRepo.UpdateOTP(email, otpCode, expiresAt); errUpdate != nil {
		return errUpdate
	}

	go func(target, name, code string) {
		if errSend := utils.SendEmailOTP(target, name, code); errSend != nil {
			log.Printf("Gagal mengirim OTP ke %s: %v\n", target, errSend)
		} else {
			log.Printf("OTP berhasil dikirim ke email: %s\n", target)
		}
	}(email, targetName, otpCode)

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

func (u *userUsecase) DeleteAccount(id string) error {
	return u.userRepo.Delete(id)
}

func (u *userUsecase) GetImpact(id string) (ConsumerImpactDTO, error) {
	repairs, co2, err := u.userRepo.GetConsumerImpact(id)
	if err != nil {
		return ConsumerImpactDTO{}, errors.New("gagal mengambil metrik dampak lingkungan")
	}

	return ConsumerImpactDTO{
		TotalRepairs: repairs,
		CO2AvoidedKg: co2,
	}, nil
}

func (u *userUsecase) UpdateProfile(id string, req UpdateProfileRequest) error {
	if strings.TrimSpace(req.FullName) == "" {
		return errors.New("nama lengkap tidak boleh kosong")
	}

	return u.userRepo.UpdateProfile(id, req.FullName, req.ProfilePictureURL)
}
