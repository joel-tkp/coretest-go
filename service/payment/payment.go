package payment

import (
	// "context"
	"encoding/json" // JSON encode
	"coretest-go/service/helper" // Helper Function
	"coretest-go/service/redis" // Redis Service Provider
	// "coretest-go/pkg/log" // Logger Package
)

// Model Definition
type Payment struct {
	ID             int64 `db:"id"`
	OrderID        int64 `db:"order_id"`
	IsConfirmed	   bool `db:"is_confirmed"`
	PaymentChannel string `db:"payment_channel"`
	Amount		   float64 `db:"amount"`
	IdempotencyKey string `db:"idempotency_key"`
}

// Constants
const SERVICEID = "payment"

// Service Implementation ------------------------------------------------------------------------------------------------
// Confirmid any payment data
func (s *Service) ConfirmPayment(paymentid int64) (error) {
	p, err := s.Get(paymentid)
	p = s.Update(p.ID, p.OrderID, true, p.PaymentChannel, p.Amount, p.IdempotencyKey)
	return err
}

// Create New Data
func (s *Service) Create(orderId int64, isConfirmed bool, paymentChannel string, amount float64, idempotencyKey string) (Payment) {
	u := Payment{
		OrderID: orderId, 
		IsConfirmed: isConfirmed, 
		PaymentChannel: paymentChannel, 
		Amount: amount,
		IdempotencyKey: idempotencyKey,
	}
	u = s.resource.Create(u)
	// invalidate cache
	helper.RevalidateCachePattern(SERVICEID + "-list-*", SERVICEID)
 	// return
    return u
}

// Get Data by ID
func (s *Service) Get(id int64) (Payment, error) {
	cacheKey := SERVICEID + "-detail-" + string(id)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse Payment;
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
func (s *Service) AllList(orderBy string, orderDirection string) ([]Payment, error) {
	cacheKey := SERVICEID + "-list-" + orderBy + orderDirection
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []Payment;
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
func (s *Service) PaginatedList(perPage int32, page int32, orderBy string, orderDirection string) ([]Payment, error) {
	cacheKey := SERVICEID + "-list-" + orderBy + orderDirection + string(perPage) + string(page)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []Payment;
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
func (s *Service) Update(id int64, orderId int64, isConfirmed bool, paymentChannel string, amount float64, idempotencyKey string) (Payment) {
	u := Payment{
		ID: id,
		OrderID: orderId, 
		IsConfirmed: isConfirmed, 
		PaymentChannel: paymentChannel, 
		Amount: amount,
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
