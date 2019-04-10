package payment

import (
	//"context"

	"encoding/json" // JSON encode
	"References/coretest/service/redis"
	"References/coretest/pkg/log"
)

type Payment struct {
	ID             int64 `db:"id"`
	OrderID        int64 `db:"order_id"`
	IsConfirmed	   bool `db:"is_confirmed"`
	PaymentChannel string `db:"payment_channel"`
	Amount		   float64 `db:"amount"`
	IdempotencyKey string `db:"idempotency_key"`
}


func (s *Service) ConfirmPayment(paymentid int64) (error) {
	p, err := s.GetPayment(paymentid)
	p = s.UpdatePayment(p.ID, p.OrderID, true, p.PaymentChannel, p.Amount, p.IdempotencyKey)
	return err
}

// Create New Payment
func (s *Service) CreatePayment(orderId int64, isConfirmed bool, paymentChannel string, amount float64, idempotencyKey string) (Payment) {
	u := s.resource.CreatePayment(orderId, isConfirmed, paymentChannel, amount, idempotencyKey)
	// begin invalidate cache
	relatedKeys,_ := redis.GetKeys("payment-list-*")
	for _, key := range relatedKeys {
		redis.Delete(key)
	}
	log.SetOutputToFile("log/payment.log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", "payment-list-*")
 	// end of cache invalidation
    return u
}

// Get Payment by ID
func (s *Service) GetPayment(/*ctx context.Context, */id int64) (Payment, error) {
	cacheKey := "payment-detail-" + string(id)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse Payment;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u, err := s.resource.GetPayment(id)
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	return u, err
}

// Get List of All Payment
func (s *Service) GetAllPayment(/*ctx context.Context, */orderBy string, orderDirection string) ([]Payment, error) {
	cacheKey := "payment-list-" + orderBy + orderDirection
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []Payment;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u,err := s.resource.GetListPayment(false, -1, -1, orderBy, orderDirection, "")
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	return u, err 
}

// Get Payment as Paginated List
func (s *Service) GetPaginatedPayment(/*ctx context.Context, */perPage int32, page int32, orderBy string, orderDirection string) ([]Payment, error) {
	cacheKey := "payment-list-" + orderBy + orderDirection + string(perPage) + string(page)
	data,_ := redis.Get(cacheKey)
	if (data != nil) {
		var cachedResponse []Payment;
		json.Unmarshal(data, &cachedResponse)
		return cachedResponse,nil
	}
	u, err := s.resource.GetListPayment(true, perPage, page, orderBy, orderDirection, "")
	b,_ := json.Marshal(u)
	redis.Set(cacheKey, b)
	return u, err
}

// Update Payment
func (s *Service) UpdatePayment(id int64, orderId int64, isConfirmed bool, paymentChannel string, amount float64, idempotencyKey string) (Payment) {
	u := s.resource.UpdatePayment(id, orderId, isConfirmed, paymentChannel, amount, idempotencyKey)
	// begin invalidate cache
	redis.Delete("payment-detail-" + string(id))
	relatedKeys,_ := redis.GetKeys("payment-list-*")
	for _, key := range relatedKeys {
		redis.Delete(key)
	}
	log.SetOutputToFile("log/payment.log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", "payment-detail-" + string(id))
	log.Infow("Cache Invalidation", "key", "payment-list-*")
 	// end of cache invalidation
    return u
}

// Delete Payment
func (s *Service) DeletePayment(/*ctx context.Context, */id int64) {
	s.resource.DeletePayment(id)
	// begin invalidate cache
	redis.Delete("payment-detail-" + string(id))
	relatedKeys,_ := redis.GetKeys("payment-list-*")
	for _, key := range relatedKeys {
		redis.Delete(key)
	}
	log.SetOutputToFile("log/payment.log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", "payment-detail-" + string(id))
	log.Infow("Cache Invalidation", "key", "payment-list-*")
 	// end of cache invalidation
}
