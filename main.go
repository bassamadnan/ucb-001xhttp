package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)



func main() {
  // Creates a router without any middleware by default
  router := gin.New()


  // Global middleware
  // Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
  // By default gin.DefaultWriter = os.Stdout
  router.Use(gin.Logger())


  // Recovery middleware recovers from any panics and writes a 500 if there was one.
  router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
    if err, ok := recovered.(string); ok {
      c.String(http.StatusInternalServerError, fmt.Sprintf("Error : %s", err))
    }
    c.AbortWithStatus(http.StatusInternalServerError)
  }))


  router.GET("/panic", func(c *gin.Context) {
    // panic with a string -- the custom middleware could save this to a database or report it to the user
    panic("500 Internal Server Error")
  })


  router.GET("/", func(c *gin.Context) {
    c.String(http.StatusOK, "200 OK")
  })


  // Listen and serve on 0.0.0.0:8080
  router.Run(":8080")
}
