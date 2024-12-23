package handlers

import (
	"fmt"

	model "github.com/bassamadnan/ucb-001xhttp/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// professor handlers
// Should be able to
// 1. View his own slots
// 2. Change slot availability
// 3. Cancel appointments

// get all slots (both available and booked)
func (h *Handler) GetAllSlots(c *gin.Context) {
	session := sessions.Default(c)
	profEmail := session.Get("email").(string)

	slots, err := model.GetProfessorSlots(h.db, profEmail)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"slots": slots})
}

// update slot availability
func (h *Handler) UpdateSlotAvailability(c *gin.Context) {
	session := sessions.Default(c)
	profEmail := session.Get("email").(string)

	var req UpdateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("error: %v\n", err)
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	err := model.UpdateSlotAvailability(h.db, profEmail, req.StartTime, *req.Available)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "slot updated successfully"})
}

// cancel appointment
func (h *Handler) CancelAppointment(c *gin.Context) {
	session := sessions.Default(c)
	userEmail := session.Get("email").(string)
	userType := session.Get("userType").(uint8)

	var req CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	err := model.CancelAppointment(h.db, userEmail, req.StartTime, userType, req.ProfessorEmail)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "appointment cancelled successfully"})
}
