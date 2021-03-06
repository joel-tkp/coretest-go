package api

import (
	//"context"

	"github.com/joel-tkp/coretest-go/service/user"
	"github.com/joel-tkp/coretest-go/service/order"
	"github.com/joel-tkp/coretest-go/service/payment"
	"github.com/joel-tkp/coretest-go/service/logistic"
)

type UserService interface {
	IsUserActive(int64) (bool, error)
	Create(string, string, bool, string) (user.User)
	Get(int64) (user.User, error)
	AllList(string, string) ([]user.User, error)
	PaginatedList(int32, int32, string, string) ([]user.User, error)
	Update(int64, string, string, bool, string) (user.User)
	Delete(int64)
}

type PaymentService interface {
	ConfirmPayment(int64) error
	Create(int64, bool, string, float64, string) (payment.Payment)
	Get(int64) (payment.Payment, error)
	AllList(string, string) ([]payment.Payment, error)
	PaginatedList(int32, int32, string, string) ([]payment.Payment, error)
	Update(int64, int64, bool, string, float64, string) (payment.Payment)
	Delete(int64)
}

type OrderService interface {
	CreateOrder(order.Order) error
	ConfirmOrder(string) error
}

type LogisticService interface {
	CreateShipment(logistic.Shipment) error
	ShipmentSent(string) error
	ShipmentReceived(string) error 
}
