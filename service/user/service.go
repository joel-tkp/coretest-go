package user

// Resource Contract
type Resource interface {
	Create(User) (User)
	Get(int64) (User, error)
	GetList(bool, int32, int32, string, string, string) ([]User, error)
	Update(User) (User)
	Delete(int64)
}

// Service Definition
type Service struct {
	resource Resource
}

// Service Initialization
func New(_resource Resource) *Service {
	s := Service{
		resource: _resource,
	}
	return &s
}
