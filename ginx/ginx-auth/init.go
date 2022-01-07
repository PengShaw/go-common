package ginxauth

import (
	"github.com/PengShaw/go-common/ginx/ginx-auth/casbinx"
	"github.com/PengShaw/go-common/ginx/ginx-auth/handlers"
	jwttoken "github.com/PengShaw/go-common/ginx/ginx-auth/jwt-token"
	"github.com/PengShaw/go-common/ginx/ginx-auth/middlewares"
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WarpAuthMiddleware(router *gin.RouterGroup, PrivateKey string, JwtExpired int, JwtTokenHeaderKey string) *gin.RouterGroup {
	jwttoken.Setup(PrivateKey, JwtExpired, JwtTokenHeaderKey)
	router.POST("/login", handlers.LoginHandler)
	authBase := router.Group("")
	authBase.Use(middlewares.AuthRequired())
	auth := authBase.Group("/auth")
	{
		auth.POST("/user", handlers.CreateUser)
		auth.GET("/user", handlers.ListUser)
		auth.DELETE("/user/:id", handlers.DeleteUser)
		auth.PUT("/user/password", handlers.UserChangePassword)
		auth.PUT("/user/group", handlers.RegroupUser)

		auth.POST("/user-group", handlers.CreateUserGroup)
		auth.GET("/user-group", handlers.ListUserGroup)
		auth.DELETE("/user-group/:id", handlers.DeleteUserGroup)
		auth.PUT("/user-group/group-name", handlers.RenameUserGroup)
		auth.PUT("/user-group/comment", handlers.CommentUserGroup)

		auth.GET("/resource", handlers.ListAuthResource)

		auth.POST("/permission", handlers.CreateAuthPermission)
		auth.GET("/permission", handlers.ListAuthPermission)
		auth.DELETE("/permission/:id", handlers.DeleteAuthPermission)

	}
	return authBase
}

func Setup(router *gin.Engine, DB *gorm.DB, exceptResource ...string) error {
	err := models.SetupDB(DB, router, exceptResource)
	if err != nil {
		return err
	}
	return casbinx.Setup()
}
