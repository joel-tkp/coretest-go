package order

import (
	//"context"
	"errors"
)

type Order struct {
	ID             int64 `db:"id"`
	UserID         int64  `db:"user_id"`
	OrderNumber    string `db:"order_number"`
	IsConfirmed	   bool   `db:"is_confirmed"`
	IdempotencyKey string `db:"idempotency_key"`
	Metadata       string `db:"metadata"`
}

func (s *Service) CreateOrder(/*ctx context.Context, */order Order) error {
	if order.IdempotencyKey == "" {
		return errors.New("order idempotency key cannot be empty")
	}
	// check wether user is active or not using user service
	if userActive, err := s.userSvc.IsUserActive(/*ctx, */order.UserID); err != nil {
		return err
	} else if !userActive {
		return errors.New("cannot create order, user is not active")
	}

	return s.resource.CreateOrder(/*ctx, */order)
}

func (s *Service) ConfirmOrder(/*ctx context.Context, */orderid string) error {
	if orderid == "" {
		return errors.New("orderid is empty")
	}
	return nil
}
