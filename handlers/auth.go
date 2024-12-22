package handlers

import (
	"fmt"

	model "github.com/bassamadnan/ucb-001xhttp/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// register new user (student/professor)
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Binding error:", err)
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// validate user type
	if req.Type != model.STUDENT_TYPE && req.Type != model.PROFESSOR_TYPE {
		c.JSON(400, gin.H{"error": "Invalid user type"})
		return
	}

	err := model.RegisterUser(h.db, req.Email, req.Name, req.Password, req.Type)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Registration successful"})
}

// login and create session
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}
	user, err := model.LoginUser(h.db, req.Email, req.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": "error loginUser"})
		return
	}
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Set("userType", user.Type)
	session.Set("email", user.Email)
	session.Save()
	c.JSON(200, gin.H{
		"message": "success",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"type":  user.Type,
		},
	})
}
