package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
)

type HttpServer struct {
	router *gin.Engine
	maker  pkg.JWTMaker

	HTTPService   services.HttpInterface
	RabbitService services.RabbitInterface
	GRPCService   services.GrpcInterface
}

func NewHttpServer(maker pkg.JWTMaker) *HttpServer {
	server := &HttpServer{
		maker: maker,
	}

	server.setRoutes()

	return server
}

func (s *HttpServer) setRoutes() {
	r := gin.Default()

	auth := r.Group("/").Use(authenticationMiddleware(s.maker)) // requires access token

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs")
	}

	r.StaticFS("/swagger", statikFs)

	r.POST("/register", s.handleRegisterUser)
	r.POST("/login", s.handleLoginUser)
	auth.POST("/payments/initiate", s.handleInitiatePayment)
	auth.GET("/payments/status/:id", s.handlePaymentPolling)

	s.router = r
}

func (s *HttpServer) Start(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
