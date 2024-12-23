package api

import (
	handler "github.com/bassamadnan/ucb-001xhttp/handlers"
	"github.com/gin-gonic/gin"
)

func InitAuth(r *gin.Engine, h *handler.Handler) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}
