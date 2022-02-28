package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go-musthave-diploma-tpl/internal/config"
	"go-musthave-diploma-tpl/internal/loyalty_system"
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

	_, err := repository.NewRepository(cfg)
	if err != nil {
		panic(fmt.Sprintf("Can't create repository: %s", err.Error()))
	}

	service := loyalty_system.New(cfg.RunAddr)

	log.Println("Starting server at port 8080")

	srv := http.Server{
		Addr:    cfg.RunAddr,
		Handler: NewRouter(service),
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func NewRouter(service *loyalty_system.Loyalty) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	return r
}
