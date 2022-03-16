package repository

import (
	"context"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/models"
)

func UploadValuesToDB(ctx context.Context, g *gophermart.Gophermart) {
	user, _ := models.NewUser("sdkfsd", "dskdjfsee")
	g.UserRepo.AddUser(ctx, user)

	account := models.NewAccount(user.UID, 23.0, 0)
	g.UserAccountRepo.AddAccount(ctx, account)

	order := models.NewOrder(user.UID, "4567890000000000")
	g.UserOrderRepo.AddUserOrder(ctx, order)
}
