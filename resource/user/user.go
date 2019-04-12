package user

import (
	"fmt"
	"coretest-go/resource/helper"
	"coretest-go/service/user"
	"github.com/jmoiron/sqlx"
	// "github.com/satori/go.uuid" // for uid primary key
)

// Backend Data Resource
type Resource struct {
	masterDB   *sqlx.DB
	followerDB *sqlx.DB
}

// Constants
const TABLE = "users"
var PRIMARYKEY = []string {"id"}
var FIELDS = []string {"id", "name", "email", "is_active", "idempotency_key"}

/**
 * SCHEMA --------------------------------------
 * ID int64 `db:"id"`
 * Name string `db:"name"`
 * Email string `db:"email"`
 * IsActive bool `db:"is_active"`
 * IdempotencyKey string `db:"idempotency_key"`
 **/

// Initiate Data Resources
func New(masterDB, followerDB *sqlx.DB) *Resource {
	r := Resource{
		masterDB:   masterDB,
		followerDB: followerDB,
	}
	return &r
}

// Create New Data
func (r *Resource) Create(u user.User) (user.User) {
	// u.ID = uuid.Must(uuid.NewV4()) // if using UUID paradigm
	tx := r.masterDB.MustBegin()
	queryBuilder := "INSERT INTO " + TABLE + " ("
	queryBuilder += helper.GenerateInsertFields(FIELDS, "", PRIMARYKEY)
	queryBuilder += ") VALUES ("
	queryBuilder += helper.GenerateInsertFields(FIELDS, ":", PRIMARYKEY)
	queryBuilder += ") RETURNING id"
    rows, err := tx.NamedQuery(queryBuilder, &u)
    tx.Commit()
    // retrieved ID
    if err != nil {
    	u.ID = 0
    	return u // return
	}
    var returnedId int64
	if rows.Next() {
		rows.Scan(&returnedId)
		u.ID = returnedId
	}
	// return
    return u
}

// Get Data by ID
func (r *Resource) Get(id int64) (user.User, error) {
	u := user.User{}
	queryBuilder := "SELECT * FROM " + TABLE + " WHERE id=$1"
    err := r.masterDB.Get(&u, queryBuilder, id)
    // rows, err = db.NamedQuery(`SELECT * FROM users WHERE id=:id`, map[string]interface{}{"id": id})
	return u, err
}

// Get List of Data
func (r *Resource) GetList(isPaginated bool, perPage int32, page int32, orderBy string, orderDirection string, extraCondition string) ([]user.User, error) {
	uList := []user.User{}

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
func (r *Resource) Update(u user.User) (user.User) {
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
    tx.NamedExec("DELETE FROM " + TABLE + " WHERE id = :id", &user.User{ID: id})
    tx.Commit()
}
