# ginx

## middlewares

```go
package main

import (
	"os"
	"fmt"
	
	"github.com/PengShaw/go-common/ginx/middlewares"
	"github.com/PengShaw/go-common/logger"
	"github.com/gin-gonic/gin"
)

func main(){
	gin.DisableConsoleColor()
	middlewares.InitJwt("PrivateKey", 10, nil)
	
	router := gin.New()
	router.Use(gin.Recovery(), middlewares.Logger())

	// ping api
	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// token
	router.POST("/api/v1/token", middlewares.LoginHandler(func(username, password string) (uint, string, bool) {
		return 1, "admin", true
	}))
	// authorized
	authorized := router.Group("/api/v1")
	authorized.Use(middlewares.JwtTokenRequired())
	{
		authorized.GET("/ping-auth", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "auth",
			})
		})
	}

	runStr := fmt.Sprintf(":%d", 8080)
	logger.Info("listen http server on " + runStr)
	if err := router.Run(runStr); err != nil {
		logger.Errorf("run http server err: %s", err.Error())
		os.Exit(1)
	}
}
```