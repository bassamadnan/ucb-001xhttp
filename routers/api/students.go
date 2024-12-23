package api

import (
	handler "github.com/bassamadnan/ucb-001xhttp/handlers"
	middleware "github.com/bassamadnan/ucb-001xhttp/middleware"
	"github.com/gin-gonic/gin"
)

func InitStudent(r *gin.Engine, h *handler.Handler) {
	student := r.Group("/api/student")
	// use the middleware we just defined
	student.Use(middleware.AuthRequired())
	student.Use(middleware.StudentRequired())
	{
		student.GET("/professors", h.GetAllProfessors)
		student.GET("/professor/:email/slots", h.GetProfessorAvailableSlots)
		student.POST("/book", h.BookProfessorSlot)
		student.GET("/appointments", h.GetMyAppointments)
		student.POST("/cancel", h.CancelAppointment)
	}
}
