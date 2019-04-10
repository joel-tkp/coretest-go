package http

import (
	"context"
	"net"
	"net/http"

	"References/coretest/api"
	"References/coretest/api/http/order"
	"References/coretest/api/http/payment"
	"References/coretest/api/http/logistic"
	"References/coretest/api/http/user"
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
	// return handler
	return router
}
