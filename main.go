package main

import (
	"fmt"
	"log"

	handlers "github.com/bassamadnan/ucb-001xhttp/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Print("1\n")
	db, err := gorm.Open(sqlite.Open("models/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("appointment-session", store))
	h := handlers.NewHandler(db)
	fmt.Print("2\n")
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
	router.Run(":8080")

}
