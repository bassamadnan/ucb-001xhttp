package main

import (
	"log"

	router "github.com/bassamadnan/ucb-001xhttp/routers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("models/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	router := router.InitRouter(db)
	router.Run(":8080")

}
