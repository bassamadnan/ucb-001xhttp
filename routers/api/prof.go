package api

import (
	handler "github.com/bassamadnan/ucb-001xhttp/handlers"
	middleware "github.com/bassamadnan/ucb-001xhttp/middleware"
	"github.com/gin-gonic/gin"
)

func InitProfessor(r *gin.Engine, h *handler.Handler) {
	professor := r.Group("/api/professor")
	// use the middleware we just defined
	professor.Use(middleware.AuthRequired())
	professor.Use(middleware.ProfessorRequired())
	{
		professor.GET("/slots", h.GetAllSlots)
		professor.PUT("/slot", h.UpdateSlotAvailability)
		professor.POST("/cancel", h.CancelAppointment)
	}
}
