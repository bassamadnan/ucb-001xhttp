package handlers

import (
	model "github.com/bassamadnan/ucb-001xhttp/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// student handlers
// students should be able to
// 1. Get available slots for a prof
// 2. To enable above we need to first get list of profs (their name and mail)
// 3. Book appointments
// 4. View self appointments

func (h *Handler) GetProfessorAvailableSlots(c *gin.Context) {
	// view available slots for a specific professor
	profEmail := c.Param("email")
	slots, err := model.GetAvailableSlots(h.db, profEmail)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"slots": slots})
}

func (h *Handler) GetAllProfessors(c *gin.Context) {
	// return the list of all the proffes (with their mail and name)
	professors, err := model.GetAllProfessors(h.db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var response []ProfessorInfo
	for _, prof := range professors {
		response = append(response, ProfessorInfo{
			Name:  prof.Name,
			Email: prof.Email,
		})
	}

	c.JSON(200, gin.H{"professors": response})
}

func (h *Handler) BookProfessorSlot(c *gin.Context) {
	// book an appointment with a professor
	session := sessions.Default(c)
	studentEmail := session.Get("email").(string)

	var req BookSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request format"})
		return
	}

	err := model.BookAppointment(h.db, studentEmail, req.StartTime, req.ProfessorEmail)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "appointment booked successfully"})
}

func (h *Handler) GetMyAppointments(c *gin.Context) {
	// view student's appointments
	session := sessions.Default(c)
	studentEmail := session.Get("email").(string)

	appointments, err := model.GetStudentAppointments(h.db, studentEmail)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"appointments": appointments})
}
