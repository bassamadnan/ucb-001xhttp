package handlers

import "gorm.io/gorm"

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	// constructor, best practice apparently.
	return &Handler{db: db}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Type     uint8  `json:"type"` // change later?
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
