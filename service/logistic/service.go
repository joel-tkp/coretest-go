package logistic

//import "context"

type Resource interface {
	CreateShipment(/*context.Context, */Shipment) error
}

type UserService interface {
	IsUserActive(/*context.Context, */int64) (bool, error)
}

type Service struct {
	resource Resource
	userSvc  UserService
}

// New payment service
func New(shipmentResource Resource, userService UserService) *Service {
	s := Service{
		resource: shipmentResource,
		userSvc:  userService,
	}
	return &s
}
