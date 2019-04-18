package logistic

import (
	// "context"

	logisticservice "github.com/joel-tkp/coretest-go/service/logistic"
	"github.com/jmoiron/sqlx"
)

type Resource struct {
	masterDB   *sqlx.DB
	followerDB *sqlx.DB
}

func New(masterDB, followerDB *sqlx.DB) *Resource {
	r := Resource{
		masterDB:   masterDB,
		followerDB: followerDB,
	}
	return &r
}

func (r *Resource) CreateShipment(/* ctx context.Context,*/ shp logisticservice.Shipment) error {
	return nil
}
