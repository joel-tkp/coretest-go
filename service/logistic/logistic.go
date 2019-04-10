package logistic

import (
	//"context"
	"errors"
)

type Shipment struct {
	ID             int64 `db:"id"`
	OrderID        int64 `db:"order_id"`
	IsSent	   	   bool `db:"is_sent"` // simple tracking (direct)
	IsReceived	   bool `db:"is_received"` // simple tracking (direct)
	Courier 	   string `db:"courier"` // simple cost (direct)
	CourierService string `db:"courier_service"` // simple cost (direct)
	Cost		   float64 `db:"cost"` // simple cost (direct)
	FromAddress	   string `db:"from_address"` // simple fullfillment (direct)
	ToAddress	   string `db:"to_address"` // simple fullfillment (direct)
	IdempotencyKey string `db:"idempotency_key"`
}

func (s *Service) CreateShipment(/*ctx context.Context, */shipment Shipment) error {
	if shipment.IdempotencyKey == "" {
		return errors.New("shipment idempotency key cannot be empty")
	}
	return nil
}

func (s *Service) ShipmentSent(/*ctx context.Context, */shipmentid string) error {
	if shipmentid == "" {
		return errors.New("shipment id is empty")
	}
	return nil
}

func (s *Service) ShipmentReceived(/*ctx context.Context, */shipmentid string) error {
	if shipmentid == "" {
		return errors.New("shipment id is empty")
	}
	return nil
}
