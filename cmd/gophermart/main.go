package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go-musthave-diploma-tpl/internal/config"
	"go-musthave-diploma-tpl/internal/controller"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg := config.LoadServerConfiguration()
	_, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	repo, err := repository.NewRepository(cfg)
	if err != nil {
		panic(fmt.Sprintf("Can't create repository: %s", err.Error()))
	}

	service := gophermart.New(cfg.RunAddr, repo)

	log.Println("Starting server at port 8080")

	srv := http.Server{
		Addr:    cfg.RunAddr,
		Handler: NewRouter(service),
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func NewRouter(service *gophermart.Loyalty) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/user/register", controller.RegisterUser(service))
	r.Post("/api/user/login", controller.LoginUser(service))
	r.Post("/api/user/orders", controller.UploadUserOrder(service))
	r.Get("/api/user/orders", controller.GetUserOrders(service))
	r.Get("/api/user/balance", controller.GetUserBalance(service))
	r.Post("/api/user/balance/withdraw", controller.WithdrawBalance(service))
	r.Get("/ping", controller.CheckConn(service))
	return r
}
