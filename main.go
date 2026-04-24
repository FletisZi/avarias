package main

import (
	"camsystem/db"
	"camsystem/router"
	"camsystem/schemas"
	"camsystem/service"
)

func main() {

	db.InitDB()

	manager := schemas.NewStreamManager()

	go service.Capturarador(manager)

	router.Initialize(manager)

}



