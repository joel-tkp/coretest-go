package payment

import (
	"fmt"
	"github.com/joel-tkp/coretest-go/resource/helper"
	"github.com/joel-tkp/coretest-go/service/payment"
	"github.com/jmoiron/sqlx"
	// "github.com/satori/go.uuid" // for uid primary key
)

// Backend Data Resource 
type Resource struct {
	masterDB   *sqlx.DB
	followerDB *sqlx.DB
}

// Constants
const TABLE = "payments"
var PRIMARYKEY = []string {"id"}
var FIELDS = []string {"id", "order_id", "is_confirmed", "payment_channel", "amount", "idempotency_key"}

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

// Create New Data
func (r *Resource) Create(u payment.Payment) (payment.Payment) {
	// u.ID = uuid.Must(uuid.NewV4()) // if using UUID paradigm
	tx := r.masterDB.MustBegin()
	queryBuilder := "INSERT INTO " + TABLE + " ("
	queryBuilder += helper.GenerateInsertFields(FIELDS, "", PRIMARYKEY)
	queryBuilder += ") VALUES ("
	queryBuilder += helper.GenerateInsertFields(FIELDS, ":", PRIMARYKEY)
	queryBuilder += ")"
    tx.NamedExec(queryBuilder, &u)
    tx.Commit()
    return u
}

// Get Data by ID
func (r *Resource) Get(id int64) (payment.Payment, error) {
	u := payment.Payment{}
	queryBuilder := "SELECT * FROM " + TABLE + " WHERE id=$1"
    err := r.masterDB.Get(&u, queryBuilder, id)
    // rows, err = db.NamedQuery(`SELECT * FROM payments WHERE id=:id`, map[string]interface{}{"id": id})
	return u, err
}

// Get List of Data
func (r *Resource) GetList(isPaginated bool, perPage int32, page int32, orderBy string, orderDirection string, extraCondition string) ([]payment.Payment, error) {
	uList := []payment.Payment{}

	_whereClause := (map[bool]string{true: " WHERE " + extraCondition + " ", false: ""})[extraCondition != ""]
	_orderDirection := (map[bool]string{true: " " + orderDirection + " ", false: " ASC "})[orderDirection != ""]
	_orderClause := (map[bool]string{true: " ORDER BY " + orderBy + _orderDirection, false: ""})[orderBy != ""]
	_perPage := (map[bool]int32{true: perPage, false: 10})[perPage > 0]
	_pageOffset := (map[bool]int32{true: ((page - 1) * _perPage), false: 0})[page > 0]
	_pagination := (map[bool]string{true: fmt.Sprintf(" LIMIT %d,%d " , _pageOffset , _perPage), false: ""})[isPaginated]

	var query string = "SELECT * FROM " + TABLE + _whereClause + _orderClause + _pagination
	err := r.masterDB.Select(&uList, query)
	return uList, err
}

// Update Data
func (r *Resource) Update(u payment.Payment) (payment.Payment) {
	tx := r.masterDB.MustBegin()
	queryBuilder := "UPDATE " + TABLE + " SET "
	queryBuilder += helper.GenerateUpdateFields(FIELDS, PRIMARYKEY)
	queryBuilder += " WHERE " 
	queryBuilder += helper.GenerateUpdateFieldsPrimaryKeyConstraint(PRIMARYKEY)
    tx.NamedExec(queryBuilder, &u)
    tx.Commit()
    return u
}

// Delete Data
func (r *Resource) Delete(id int64) {
	tx := r.masterDB.MustBegin()
    tx.NamedExec("DELETE FROM " + TABLE + " WHERE id = :id", &payment.Payment{ID: id})
    tx.Commit()
}
