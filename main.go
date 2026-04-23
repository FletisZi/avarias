package main

import (

	"camsystem/router"
	"camsystem/db"
)


func main() {

	db.InitDB()

	router.Initialize()

}