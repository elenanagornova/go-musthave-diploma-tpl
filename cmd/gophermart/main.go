package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/config"
	"go-musthave-diploma-tpl/internal/controller"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/repository"
	"go-musthave-diploma-tpl/internal/repository/account"
	"go-musthave-diploma-tpl/internal/repository/user"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	cfg := config.LoadServerConfiguration()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	m, err := repository.RunMigration(cfg.DatabaseDSN)
	if err != nil && !m {
		log.Fatal(err)
	}

	userRepo := user.NewUserRepository(conn)
	userAccountRepo := account.NewUserAccountRepository(conn)
	service := gophermart.NewGophermart(cfg.RunAddr, userRepo, userAccountRepo)

	log.Println("Starting server at port 8080")

	srv := http.Server{
		Addr:    cfg.RunAddr,
		Handler: controller.NewRouter(ctx, service),
	}

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
