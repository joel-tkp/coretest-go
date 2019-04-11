package payment

import (
	"net/http"
    "github.com/julienschmidt/httprouter"
	"strconv" // String Converter
	"time" // time lib

	"References/coretest/api"
	"References/coretest/api/http/base"
	"References/coretest/service/payment"
)

var paymentService api.PaymentService

func Init(service api.PaymentService) {
	paymentService = service
}

func ConfirmPayment(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	objId,_ := strconv.ParseInt(qs.ByName("id"), 0, 64)
	u, err := paymentService.Get(objId)
	if err != nil {
		base.WriteSuccessResponse(w, base.Error, "Payment not found!", start, u)
	}
	u = paymentService.Update(u.ID, u.OrderID, true, u.PaymentChannel, u.Amount, u.IdempotencyKey)
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "Payment for order #" + string(u.OrderID) + " confirmed!", start, u)
}

func Create(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	orderId,_ := strconv.ParseInt(r.FormValue("orderId"), 0, 64)
	isConfirmed,_ := strconv.ParseBool(r.FormValue("isConfirmed"))
	amount,_ := strconv.ParseFloat(r.FormValue("amount"), 64)
	u := paymentService.Create(orderId, isConfirmed, r.FormValue("paymentChannel"), amount, "testIdempotencyKey")
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "Payment for order #" + string(u.OrderID) + " created!", start, u)
}

func Detail(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	objId,_ := strconv.ParseInt(qs.ByName("id"), 0, 64)
	u, err := paymentService.Get(objId)
	if err != nil {
		base.WriteSuccessResponse(w, base.Error, "Payment not found!", start, u)
	}
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "Payment for order #" + string(u.OrderID) + " retrieved!", start, u)
}

func List(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	orderBy := r.URL.Query().Get("orderBy")
	orderDir := r.URL.Query().Get("orderDir")
	paginated := r.URL.Query().Get("paginated")

	var objectList []payment.Payment
	if paginated == "no" {
		objectList, _ = paymentService.AllList(orderBy, orderDir)
	} else {
		page,_ := strconv.ParseInt(r.URL.Query().Get("page"), 0, 64)
		objectList, _ = paymentService.PaginatedList(10, int32(page), orderBy, orderDir)
	}
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "Payment list retrieved!", start, objectList)
}

func Update(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	objId,_ := strconv.ParseInt(qs.ByName("id"), 0, 64)
	u, err := paymentService.Get(objId)
	if err != nil {
		base.WriteSuccessResponse(w, base.Error, "Payment not found!", start, u)
	}
	orderId,_ := strconv.ParseInt(r.FormValue("orderId"), 0, 64)
	isConfirmed,_ := strconv.ParseBool(r.FormValue("isConfirmed"))
	amount,_ := strconv.ParseFloat(r.FormValue("amount"), 64)
	u = paymentService.Update(u.ID, orderId, isConfirmed, r.FormValue("paymentChannel"), amount, "testIdempotencyKey")
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "Payment for order #" + string(u.OrderID) + " updated!", start, u)
}

func Delete(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	start := time.Now()
	objId,_ := strconv.ParseInt(qs.ByName("id"), 0, 64)
	u, err := paymentService.Get(objId)
	if err != nil {
		base.WriteSuccessResponse(w, base.Error, "Payment not found!", start, u)
	}
	paymentService.Delete(u.ID)
	// return response as JSON
	base.WriteSuccessResponse(w, base.Success, "Payment for order #" + string(u.OrderID) + " created!", start, u)
}
