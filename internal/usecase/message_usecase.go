package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/repository"
)

type messageUsecase struct {
	messageRepo domain.MessageRepository
	orderRepo   repository.OrderRepository
}

func NewMessageUsecase(messageRepo domain.MessageRepository, orderRepo repository.OrderRepository) domain.MessageUsecase {
	return &messageUsecase{
		messageRepo: messageRepo,
		orderRepo:   orderRepo,
	}
}

func (u *messageUsecase) SendMessage(orderID, senderID, senderRole, content string) (*domain.Message, error) {
	parsedOrderID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.New("invalid order id format")
	}

	parsedSenderID, err := uuid.Parse(senderID)
	if err != nil {
		return nil, errors.New("invalid sender id format")
	}

	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}

	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if senderRole == "user" && order.CustomerID.String() != senderID {
		return nil, errors.New("unauthorized access to order")
	}
	if senderRole == "technician" && (order.TechnicianID == nil || order.TechnicianID.String() != senderID) {
		return nil, errors.New("unauthorized access to order")
	}

	message := &domain.Message{
		OrderID:    parsedOrderID,
		SenderID:   parsedSenderID,
		SenderRole: senderRole,
		Content:    content,
	}

	if err := u.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (u *messageUsecase) GetOrderMessages(orderID string) ([]domain.Message, error) {
	if _, err := uuid.Parse(orderID); err != nil {
		return nil, errors.New("invalid order id format")
	}
	return u.messageRepo.GetByOrderID(orderID)
}
