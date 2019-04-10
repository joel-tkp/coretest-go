package user

import (
	// "context"
	"encoding/json" // JSON encode
	"References/coretest/service/redis"
	"References/coretest/pkg/log"
)

type User struct {
	ID int64 `db:"id"`
	Name string `db:"name"`
	Email string `db:"email"`
	IsActive bool `db:"is_active"`
	IdempotencyKey string `db:"idempotency_key"`
}

// Check wether an user is active or not
func (s *Service) IsUserActive(/*ctx context.Context, */userid int64) (bool, error) {
	u, err := s.GetUser(userid)
	return u.IsActive, err
}

// Create New User
func (s *Service) CreateUser(/*ctx context.Context, */name string, email string, isActive bool, idempotencyKey string) (User) {
	u := s.resource.CreateUser(name, email, isActive, idempotencyKey)
	// begin invalidate cache
	relatedKeys,_ := redis.GetKeys("user-list-*")
	for _, key := range relatedKeys {
		redis.Delete(key)
	}
	log.SetOutputToFile("log/user.log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", "user-list-*")
 	// end of cache invalidation
    return u
}

// Get User by ID
func (s *Service) GetUser(/*ctx context.Context, */id int64) (User, error) {
	cacheKey := "user-detail-" + string(id)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse User;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u, err := s.resource.GetUser(id)
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	return u, err
}

// Get List of All User
func (s *Service) GetAllUser(/*ctx context.Context, */orderBy string, orderDirection string) ([]User, error) {
	cacheKey := "user-list-" + orderBy + orderDirection
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []User;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u,err := s.resource.GetListUser(false, -1, -1, orderBy, orderDirection, "")
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	return u, err 
}

// Get User as Paginated List
func (s *Service) GetPaginatedUser(/*ctx context.Context, */perPage int32, page int32, orderBy string, orderDirection string) ([]User, error) {
	cacheKey := "user-list-" + orderBy + orderDirection + string(perPage) + string(page)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []User;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u, err := s.resource.GetListUser(true, perPage, page, orderBy, orderDirection, "")
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	return u, err
}

// Update User
func (s *Service) UpdateUser(/*ctx context.Context, */id int64, name string, email string, isActive bool, idempotencyKey string) (User) {
	u := s.resource.UpdateUser(id, name, email, isActive, idempotencyKey)
	// begin invalidate cache
	redis.Delete("user-detail-" + string(id))
	relatedKeys,_ := redis.GetKeys("user-list-*")
	for _, key := range relatedKeys {
		redis.Delete(key)
	}
	log.SetOutputToFile("log/user.log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", "user-detail-" + string(id))
	log.Infow("Cache Invalidation", "key", "user-list-*")
 	// end of cache invalidation
    return u
}

// Delete User
func (s *Service) DeleteUser(/*ctx context.Context, */id int64) {
	s.resource.DeleteUser(id)
	// begin invalidate cache
	redis.Delete("user-detail-" + string(id))
	relatedKeys,_ := redis.GetKeys("user-list-*")
	for _, key := range relatedKeys {
		redis.Delete(key)
	}
	log.SetOutputToFile("log/user.log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", "user-detail-" + string(id))
	log.Infow("Cache Invalidation", "key", "user-list-*")
 	// end of cache invalidation
}