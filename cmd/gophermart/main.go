package main

import (
	"context"
	"github.com/jackc/pgx/v4"
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

	userRepo := repository.NewUserRepository(conn)
	userAccountRepo := repository.NewUserAccountRepository(conn)
	userOrderRepo := repository.NewUserOrderRepository(conn)
	withdrawalRepo := repository.NewWithdrawalRepository(conn)
	service := gophermart.NewGophermart(cfg.RunAddr, cfg.AccuralSystemAddr, userRepo, userAccountRepo, userOrderRepo, withdrawalRepo)

	log.Println("Starting server at port 8080")

	repository.UploadValuesToDb(ctx, service)

	go service.UpdateOrders(ctx)
	srv := http.Server{
		Addr:    cfg.RunAddr,
		Handler: controller.NewRouter(ctx, service),
	}

	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
