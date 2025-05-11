package controllers

import (
	"github.com/championswimmer/api.midpoint.place/src/db"
	"gorm.io/gorm"
)

type UsersController struct {
	db *gorm.DB
}

func CreateUsersController() *UsersController {
	appDb := db.GetAppDB()
	return &UsersController{
		db: appDb,
	}
}
