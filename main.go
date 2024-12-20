package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"uniqueIndex;size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	db, _ := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	db.AutoMigrate(&User{})
	var users []User
	user := User{Name: "name1", Email: "name1@mail"}
	users = append(users, user)
	result := db.Create(&users)
	result = db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	println(result)

}
