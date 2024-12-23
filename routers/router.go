package router

import (
	handlers "github.com/bassamadnan/ucb-001xhttp/handlers"
	api "github.com/bassamadnan/ucb-001xhttp/routers/api"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("appointment-session", store))
	h := handlers.NewHandler(db)

	api.InitAuth(router, h)
	api.InitStudent(router, h)
	api.InitProfessor(router, h)
	return router
}
