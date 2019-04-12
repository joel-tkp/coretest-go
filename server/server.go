package server

import (
	"io/ioutil"
	"net"

	_ "github.com/lib/pq" // backend-db driver
	httpapi "coretest-go/api/http" // HTTP handler
	orderresource "coretest-go/resource/order" // resource entity
	paymentresource "coretest-go/resource/payment" // resource entity
	logisticresource "coretest-go/resource/logistic" // resource entity
	userresource "coretest-go/resource/user" // resource entity
	orderservice "coretest-go/service/order" // service provider
	paymentservice "coretest-go/service/payment" // service provider
	logisticservice "coretest-go/service/logistic" // service provider
	userservice "coretest-go/service/user" // service provider
	"github.com/jmoiron/sqlx" // backend-db wrapper extension
	"gopkg.in/yaml.v2" // config-read
	redis "coretest-go/service/redis"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name VARCHAR (128) NOT NULL,
    email VARCHAR (255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    idempotency_key VARCHAR (1024) NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id serial PRIMARY KEY,
    user_id INTEGER NULL,
    order_number VARCHAR (64) NOT NULL,
    is_confirmed BOOLEAN DEFAULT FALSE,
    idempotency_key VARCHAR (1024) NULL,
    metadata TEXT NULL,
  	CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id)
	 	REFERENCES users (id) MATCH SIMPLE
	  	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS payments (
	id serial PRIMARY KEY,
	order_id INTEGER NULL,
	is_confirmed BOOLEAN DEFAULT FALSE,
	payment_channel VARCHAR (64),
	amount FLOAT,
    idempotency_key VARCHAR (1024) NULL,
  	CONSTRAINT payments_order_id_fkey FOREIGN KEY (order_id)
	 	REFERENCES orders (id) MATCH SIMPLE
	  	ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS shipments (
	id serial PRIMARY KEY,
	order_id INTEGER NULL,
	is_sent BOOLEAN DEFAULT FALSE,
	is_received BOOLEAN DEFAULT FALSE,
	courier VARCHAR (128),
	courier_service VARCHAR (128),
	cost FLOAT,
	from_address VARCHAR (1024),
	to_address VARCHAR (1024),
    idempotency_key VARCHAR (1024) NULL,
  	CONSTRAINT payments_order_id_fkey FOREIGN KEY (order_id)
	 	REFERENCES orders (id) MATCH SIMPLE
	  	ON UPDATE NO ACTION ON DELETE NO ACTION
);`

// Main program or run the server
func Main() error {
	// read config from config directory
	out, err := ioutil.ReadFile("config/coretest.config.yml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(out, &Config); err != nil {
		return err
	}

	masterDB, err := sqlx.Connect("postgres", Config.Database.Master)
	if err != nil {
		return err
	}
	followerDB, err := sqlx.Connect("postgres", Config.Database.Follower)
	if err != nil {
		return err
	}

	// Migration First
	masterDB.MustExec(schema)
	followerDB.MustExec(schema)

	// Init Redis Cache Service
	redis.InitService(Config.Redis.Address)

	// user
	userres := userresource.New(masterDB, followerDB)
	usersvc := userservice.New(userres)
	// order
	orderres := orderresource.New(masterDB, followerDB)
	ordersvc := orderservice.New(orderres, usersvc)
	// payment
	paymentres := paymentresource.New(masterDB, followerDB)
	paymentsvc := paymentservice.New(paymentres, usersvc)
	// shipment
	shipmentres := logisticresource.New(masterDB, followerDB)
	shipmentsvc := logisticservice.New(shipmentres, usersvc)

	// create a new listener for http and grpc server
	listener, err := net.Listen("tcp", Config.Server.Host + ":" + Config.Server.Port)
	if err != nil {
		return err
	}

	httpserver := httpapi.Server{
		UserService:    usersvc,
		OrderService:   ordersvc,
		PaymentService: paymentsvc,
		LogisticService: shipmentsvc}

	return httpserver.Serve(listener)
}
