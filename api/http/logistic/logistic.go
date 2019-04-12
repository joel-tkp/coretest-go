package logistic

import (
	"net/http"

	"coretest-go/api"
)

var logisticService api.LogisticService

func Init(service api.LogisticService) {
	logisticService = service
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {

}

func ConfirmOrder(w http.ResponseWriter, r *http.Request) {

}
