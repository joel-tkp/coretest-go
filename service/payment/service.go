package payment

// import "context"

type Resource interface {
	Create(Payment) (Payment)
	Get(int64) (Payment, error)
	GetList(bool, int32, int32, string, string, string) ([]Payment, error)
	Update(Payment) (Payment)
	Delete(int64)
}

type UserService interface {
	IsUserActive(int64) (bool, error)
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
