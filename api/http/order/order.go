package order

import (
	"net/http"

	"coretest-go/api"
)

var orderService api.OrderService

func Init(service api.OrderService) {
	orderService = service
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {

}

func ConfirmOrder(w http.ResponseWriter, r *http.Request) {

}
