package payment

import (
	// "context"
	"fmt"

	"References/coretest/service/payment"
	"github.com/jmoiron/sqlx"
)

type Resource struct {
	masterDB   *sqlx.DB
	followerDB *sqlx.DB
}

/**
 * SCHEMA --------------------------------------
 * ID             int64 `db:"id"`
 * OrderID        int64 `db:"order_id"`
 * IsConfirmed	   bool `db:"is_confirmed"`
 * PaymentChannel string `db:"payment_channel"`
 * Amount		   float64 `db:"amount"`
 * IdempotencyKey string `db:"idempotency_key"`
 **/

func New(masterDB, followerDB *sqlx.DB) *Resource {
	r := Resource{
		masterDB:   masterDB,
		followerDB: followerDB,
	}
	return &r
}


// Create New Payment
func (r *Resource) CreatePayment(orderId int64, isConfirmed bool, paymentChannel string, amount float64, idempotencyKey string) (payment.Payment) {
	// uid := uuid.Must(uuid.NewV4())
	p := payment.Payment{
		// ID: uid,
		OrderID: orderId, 
		IsConfirmed: isConfirmed, 
		PaymentChannel: paymentChannel, 
		Amount: amount,
		IdempotencyKey: idempotencyKey,
	}
	tx := r.masterDB.MustBegin()
    tx.NamedExec("INSERT INTO payments (order_id, is_confirmed, payment_channel, amount, idempotency_key) VALUES (:order_id, :is_confirmed, :payment_channel, :amount, :idempotency_key)", &p)
    tx.Commit()
    return p
}

// Get Payment by ID
func (r *Resource) GetPayment(id int64) (payment.Payment, error) {
	u := payment.Payment{}
    err := r.masterDB.Get(&u, "SELECT * FROM payments WHERE id=$1", id)
    // rows, err = db.NamedQuery(`SELECT * FROM payments WHERE id=:id`, map[string]interface{}{"id": id})
	return u, err
}

// Get List of Payment
func (r *Resource) GetListPayment(isPaginated bool, perPage int32, page int32, orderBy string, orderDirection string, extraCondition string) ([]payment.Payment, error) {
	uList := []payment.Payment{}

	_whereClause := (map[bool]string{true: " WHERE " + extraCondition + " ", false: ""})[extraCondition != ""]
	_orderDirection := (map[bool]string{true: " " + orderDirection + " ", false: " ASC "})[orderDirection != ""]
	_orderClause := (map[bool]string{true: " ORDER BY " + orderBy + _orderDirection, false: ""})[orderBy != ""]
	_perPage := (map[bool]int32{true: perPage, false: 10})[perPage > 0]
	_pageOffset := (map[bool]int32{true: ((page - 1) * _perPage), false: 0})[page > 0]
	_pagination := (map[bool]string{true: fmt.Sprintf(" LIMIT %d,%d " , _pageOffset , _perPage), false: ""})[isPaginated]

	var query string = "SELECT * FROM payments" + _whereClause + _orderClause + _pagination
	err := r.masterDB.Select(&uList, query)
	return uList, err
}

// Update Payment
func (r *Resource) UpdatePayment(id int64, orderId int64, isConfirmed bool, paymentChannel string, amount float64, idempotencyKey string) (payment.Payment) {
	p := payment.Payment{
		ID: id,
		OrderID: orderId, 
		IsConfirmed: isConfirmed, 
		PaymentChannel: paymentChannel, 
		Amount: amount,
		IdempotencyKey: idempotencyKey,
	}
	tx := r.masterDB.MustBegin()
    tx.NamedExec("UPDATE payments SET order_id = :order_id, is_confirmed = :is_confirmed, payment_channel = :payment_channel, amount = :amount, idempotency_key = :idempotency_key WHERE id = :id", &p)
    tx.Commit()
    return p
}

// Delete Payment
func (r *Resource) DeletePayment(id int64) {
	tx := r.masterDB.MustBegin()
    tx.NamedExec("DELETE FROM payments WHERE id = :id", &payment.Payment{ID: id})
    tx.Commit()
}

/*
func (r *Resource) CreatePayment(pym payment.Payment) error {
	return nil
}
*/
