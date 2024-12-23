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

type ProfessorInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type BookSlotRequest struct {
	ProfessorEmail string `json:"professor_email" binding:"required,email"`
	StartTime      uint8  `json:"start_time" binding:"required"`
}

type UpdateSlotRequest struct {
	StartTime uint8 `json:"start_time" binding:"required"`
	Available *bool `json:"available" binding:"required"`
}

type CancelAppointmentRequest struct {
	StartTime      uint8  `json:"start_time" binding:"required"`
	ProfessorEmail string `json:"professor_email,omitempty"`
}
