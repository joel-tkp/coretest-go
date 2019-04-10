package api

import (
	//"context"

	"References/coretest/service/user"
	"References/coretest/service/order"
	"References/coretest/service/payment"
	"References/coretest/service/logistic"
)

type UserService interface {
	IsUserActive(/*context.Context, */int64) (bool, error)
	CreateUser(/*context.Context, */string, string, bool, string) (user.User)
	GetUser(/*context.Context, */int64) (user.User, error)
	GetAllUser(/*context.Context, */string, string) ([]user.User, error)
	GetPaginatedUser(/*context.Context, */int32, int32, string, string) ([]user.User, error)
	UpdateUser(/*context.Context, */int64, string, string, bool, string) (user.User)
	DeleteUser(/*context.Context, */int64)
}

type PaymentService interface {
	ConfirmPayment(/*context.Context, */int64) error
	CreatePayment(/*context.Context, */int64, bool, string, float64, string) (payment.Payment)
	GetPayment(/*context.Context, */int64) (payment.Payment, error)
	GetAllPayment(/*context.Context, */string, string) ([]payment.Payment, error)
	GetPaginatedPayment(/*context.Context, */int32, int32, string, string) ([]payment.Payment, error)
	UpdatePayment(/*context.Context, */int64, int64, bool, string, float64, string) (payment.Payment)
	DeletePayment(/*context.Context, */int64)
}

type OrderService interface {
	CreateOrder(/*context.Context, */order.Order) error
	ConfirmOrder(/*context.Context, */string) error
}

type LogisticService interface {
	CreateShipment(/*context.Context, */logistic.Shipment) error
	ShipmentSent(/*context.Context, */string) error
	ShipmentReceived(/*context.Context, */string) error 
}
