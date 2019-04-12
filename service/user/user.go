package user

import (
	// "context"
	"encoding/json" // JSON encode
	"coretest-go/service/helper" // Helper Function
	"coretest-go/service/redis" // Redis Service Provider
	// "coretest-go/pkg/log" // Logger Package
)

// Model Definition
type User struct {
	ID int64 `db:"id"`
	Name string `db:"name"`
	Email string `db:"email"`
	IsActive bool `db:"is_active"`
	IdempotencyKey string `db:"idempotency_key"`
}

// Constants
const SERVICEID = "user"

// Service Implementation ------------------------------------------------------------------------------------------------
// Check wether an user is active or not
func (s *Service) IsUserActive(userid int64) (bool, error) {
	u, err := s.Get(userid)
	return u.IsActive, err
}

// Create New Data
func (s *Service) Create(name string, email string, isActive bool, idempotencyKey string) (User) {
	u := User{
		Name: name, 
		Email: email, 
		IsActive: isActive, 
		IdempotencyKey: idempotencyKey,
	}
	u = s.resource.Create(u)
	// invalidate cache
	helper.RevalidateCachePattern(SERVICEID + "-list-*", SERVICEID)
 	// return
    return u
}

// Get Data by ID
func (s *Service) Get(id int64) (User, error) {
	cacheKey := SERVICEID + "-detail-" + string(id)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse User;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u, err := s.resource.Get(id)
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	// return
	return u, err
}

// Get List of All Data
func (s *Service) AllList(orderBy string, orderDirection string) ([]User, error) {
	cacheKey := SERVICEID + "-list-" + orderBy + orderDirection
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []User;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u,err := s.resource.GetList(false, -1, -1, orderBy, orderDirection, "")
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	// return
	return u, err 
}

// Get Data as Paginated List
func (s *Service) PaginatedList(perPage int32, page int32, orderBy string, orderDirection string) ([]User, error) {
	cacheKey := SERVICEID + "-list-" + orderBy + orderDirection + string(perPage) + string(page)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []User;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u, err := s.resource.GetList(true, perPage, page, orderBy, orderDirection, "")
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	// return
	return u, err
}

// Update Data
func (s *Service) Update(id int64, name string, email string, isActive bool, idempotencyKey string) (User) {
	u := User{
		ID: id,
		Name: name, 
		Email: email, 
		IsActive: isActive, 
		IdempotencyKey: idempotencyKey,
	}
	u = s.resource.Update(u)
	// invalidate cache
	helper.RevalidateCache(SERVICEID + "-detail-" + string(id), SERVICEID)
	helper.RevalidateCachePattern(SERVICEID + "-list-*", SERVICEID)
 	// return
    return u
}

// Delete Data
func (s *Service) Delete(id int64) {
	s.resource.Delete(id)
	// invalidate cache
	helper.RevalidateCache(SERVICEID + "-detail-" + string(id), SERVICEID)
	helper.RevalidateCachePattern(SERVICEID + "-list-*", SERVICEID)
}
