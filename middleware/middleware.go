package middleware

import (
	model "github.com/bassamadnan/ucb-001xhttp/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	// check if user is authenticated, will be used ALWAYS after auth , this may not be necessary if we have
	// model.STUDENT_TYPE or PROF_TYPE but its best to follow standard practice
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userEmail := session.Get("email") // this must've been set if logged in
		if userEmail == nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func StudentRequired() gin.HandlerFunc {
	// for student related endpoints
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userType := session.Get("userType")
		if userType == nil || userType.(uint8) != model.STUDENT_TYPE {
			c.JSON(403, gin.H{"error": "forbidden: student access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ProfessorRequired() gin.HandlerFunc {
	// for prof related endpoints
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userType := session.Get("userType")
		if userType == nil || userType.(uint8) != model.PROFESSOR_TYPE {
			c.JSON(403, gin.H{"error": "forbidden: professor access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
