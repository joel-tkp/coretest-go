package user

import (
	"fmt"
	"References/coretest/service/user"
	"github.com/jmoiron/sqlx"
	// "github.com/satori/go.uuid" // for uid primary key
)

type Resource struct {
	masterDB   *sqlx.DB
	followerDB *sqlx.DB
}

/**
 * SCHEMA --------------------------------------
 * ID int64 `db:"id"`
 * Name string `db:"name"`
 * Email string `db:"email"`
 * IsActive bool `db:"is_active"`
 * IdempotencyKey string `db:"idempotency_key"`
 **/

// Initiate User Resources
func New(masterDB, followerDB *sqlx.DB) *Resource {
	r := Resource{
		masterDB:   masterDB,
		followerDB: followerDB,
	}
	return &r
}

// Create New User
func (r *Resource) CreateUser(name string, email string, isActive bool, idempotencyKey string) (user.User) {
	// uid := uuid.Must(uuid.NewV4())
	u := user.User{
		// ID: uid,
		Name: name, 
		Email: email, 
		IsActive: isActive, 
		IdempotencyKey: idempotencyKey,
	}
	tx := r.masterDB.MustBegin()
    tx.NamedExec("INSERT INTO users (name, email, is_active, idempotency_key) VALUES (:name, :email, :is_active, :idempotency_key)", &u)
    tx.Commit()
    return u
}

// Get User by ID
func (r *Resource) GetUser(id int64) (user.User, error) {
	u := user.User{}
    err := r.masterDB.Get(&u, "SELECT * FROM users WHERE id=$1", id)
    // rows, err = db.NamedQuery(`SELECT * FROM users WHERE id=:id`, map[string]interface{}{"id": id})
	return u, err
}

// Get List of User
func (r *Resource) GetListUser(isPaginated bool, perPage int32, page int32, orderBy string, orderDirection string, extraCondition string) ([]user.User, error) {
	uList := []user.User{}

	_whereClause := (map[bool]string{true: " WHERE " + extraCondition + " ", false: ""})[extraCondition != ""]
	_orderDirection := (map[bool]string{true: " " + orderDirection + " ", false: " ASC "})[orderDirection != ""]
	_orderClause := (map[bool]string{true: " ORDER BY " + orderBy + _orderDirection, false: ""})[orderBy != ""]
	_perPage := (map[bool]int32{true: perPage, false: 10})[perPage > 0]
	_pageOffset := (map[bool]int32{true: ((page - 1) * _perPage), false: 0})[page > 0]
	_pagination := (map[bool]string{true: fmt.Sprintf(" LIMIT %d,%d " , _pageOffset , _perPage), false: ""})[isPaginated]

	var query string = "SELECT * FROM users" + _whereClause + _orderClause + _pagination
	err := r.masterDB.Select(&uList, query)
	return uList, err
}

// Update User
func (r *Resource) UpdateUser(id int64, name string, email string, isActive bool, idempotencyKey string) (user.User) {
	u := user.User{
		ID: id,
		Name: name, 
		Email: email, 
		IsActive: isActive, 
		IdempotencyKey: idempotencyKey,
	}
	tx := r.masterDB.MustBegin()
    tx.NamedExec("UPDATE users SET name = :name, email = :email, is_active = :is_active, idempotency_key = :idempotency_key WHERE id = :id", &u)
    tx.Commit()
    return u
}

// Delete User
func (r *Resource) DeleteUser(id int64) {
	tx := r.masterDB.MustBegin()
    tx.NamedExec("DELETE FROM users WHERE id = :id", &user.User{ID: id})
    tx.Commit()
}
