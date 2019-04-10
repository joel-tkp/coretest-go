package user

// Contract Resource of user
type Resource interface {
	CreateUser(string, string, bool, string) (User)
	GetUser(int64) (User, error)
	GetListUser(bool, int32, int32, string, string, string) ([]User, error)
	UpdateUser(int64, string, string, bool, string) (User)
	DeleteUser(int64)
}

// Service of user
type Service struct {
	resource Resource
}

// New user service
func New(userResource Resource) *Service {
	s := Service{
		resource: userResource,
	}
	return &s
}
