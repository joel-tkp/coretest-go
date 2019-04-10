package logistic

import (
	"net/http"

	"References/coretest/api"
)

var logisticService api.LogisticService

func Init(service api.LogisticService) {
	logisticService = service
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {

}

func ConfirmOrder(w http.ResponseWriter, r *http.Request) {

}
