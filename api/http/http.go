package http

import (
	"context"
	"net"
	"net/http"

	"github.com/joel-tkp/coretest-go/api"
	"github.com/joel-tkp/coretest-go/api/http/order"
	"github.com/joel-tkp/coretest-go/api/http/payment"
	"github.com/joel-tkp/coretest-go/api/http/logistic"
	"github.com/joel-tkp/coretest-go/api/http/user"
    "github.com/julienschmidt/httprouter"
)

type Server struct {
	server         *http.Server
	UserService    api.UserService
	OrderService   api.OrderService
	PaymentService api.PaymentService
	LogisticService api.LogisticService
}

func (s *Server) Serve(lis net.Listener) error {
	s.server = &http.Server{}

	// init all handler
	user.Init(s.UserService)
	order.Init(s.OrderService)
	payment.Init(s.PaymentService)
	logistic.Init(s.LogisticService)
	// import all route into server handler
	s.server.Handler = handler()

	return s.server.Serve(lis)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func Index(w http.ResponseWriter, r *http.Request, qs httprouter.Params) {
	w.Write([]byte("Core test!"))
}

func handler() *httprouter.Router {
    router := httprouter.New()
    // base context
    router.GET("/", Index)
    // user context
    router.GET("/api/v1/user", user.List)
    router.POST("/api/v1/user/update/:id", user.Update)
    router.POST("/api/v1/user/create", user.Create)
    router.GET("/api/v1/user/:id", user.Detail)
    router.POST("/api/v1/user/delete/:id", user.Delete)
    // payment context
    router.GET("/api/v1/payment", payment.List)
    router.POST("/api/v1/payment/update/:id", payment.Update)
    router.POST("/api/v1/payment/confirm/:id", payment.ConfirmPayment)
    router.POST("/api/v1/payment/create", payment.Create)
    router.GET("/api/v1/payment/:id", payment.Detail)
    router.POST("/api/v1/payment/delete/:id", payment.Delete)
	// return handler
	return router
}
