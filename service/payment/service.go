package payment

// import "context"

type Resource interface {
	CreatePayment(int64, bool, string, float64, string) (Payment)
	GetPayment(int64) (Payment, error)
	GetListPayment(bool, int32, int32, string, string, string) ([]Payment, error)
	UpdatePayment(int64, int64, bool, string, float64, string) (Payment)
	DeletePayment(int64)
}

type UserService interface {
	IsUserActive(/*context.Context, */int64) (bool, error)
}

type Service struct {
	resource Resource
	userSvc  UserService
}

// New payment service
func New(paymentResource Resource, userService UserService) *Service {
	s := Service{
		resource: paymentResource,
		userSvc:  userService,
	}
	return &s
}
