package main

import (
	"booking_togo/internal/config"
	deliveryHttp "booking_togo/internal/delivery/http"
	"booking_togo/internal/middleware"
	"booking_togo/internal/repository"
	"booking_togo/internal/usecase"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()
	config.InitLogger()
	// connection DB with PGX
	pgxPool, pgxPoolErr := config.NewPool(cfg)
	if pgxPoolErr != nil {
		log.Fatalf("failed to connect to database: %v", pgxPoolErr)
	}

	// repository
	repo := repository.NewUserRepository(pgxPool)

	// usecase
	usecaseUser := usecase.NewUserUsecase(repo)

	// handlers
	h := deliveryHttp.NewUserFamilyHandler(usecaseUser)

	// router
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	api := r.PathPrefix("/api/v1").Subrouter()
	h.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
