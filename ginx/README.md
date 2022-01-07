# ginx

```go
package main

import (
	"fmt"
	"os"
	
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"

	ginxlogger "github.com/PengShaw/go-common/ginx/ginx-logger"
	ginxauth "github.com/PengShaw/go-common/ginx/ginx-auth"
	"github.com/PengShaw/go-common/logger"
)

var db *gorm.DB

func main(){
	gin.DisableConsoleColor()
	router := gin.New()
	// 1. ginx-logger middleware
	router.Use(gin.Recovery(), ginxlogger.Middleware())

	api := router.Group("/api")
	// 2. ginx-auth middleware
	authApi := ginxauth.WarpAuthMiddleware(api, "private-key", 6000, "JWT-Token")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	authApi.GET("/ping-auth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong-auth",
		})
	})

	// 3. ginx-auth setup resource for permission control
	err := ginxauth.Setup(router, db, "/api/login")
	if err != nil {
		panic(err.Error())
	}

	logger.Info("listen http server on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Errorf("run http server err: %s", err.Error())
		os.Exit(1)
	}
}
```