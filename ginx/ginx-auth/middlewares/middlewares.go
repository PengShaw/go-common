package middlewares

import (
	"errors"
	"github.com/PengShaw/go-common/ginx/ginx-auth/casbinx"
	jwttoken "github.com/PengShaw/go-common/ginx/ginx-auth/jwt-token"
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	contextkey "github.com/PengShaw/go-common/ginx/ginx-context-key"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(jwttoken.TokenHeaderKey)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "missing " + jwttoken.TokenHeaderKey + " in Header."})
			return
		}

		claims, err := jwttoken.Verify(tokenString)
		if err != nil {
			if errors.Is(err, jwttoken.ErrTokenFormat) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				return
			} else if errors.Is(err, jwttoken.ErrTokenExpired) || errors.Is(err, jwttoken.ErrTokenNoActive) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}
		}

		user, err := models.GetUser(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		if !casbinx.CheckPermission(user.AuthUserGroupID, c.FullPath(), c.Request.Method) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "no permission"})
			return
		}

		c.Set(contextkey.UserIDContextKey, claims.UserID)
		c.Next()
	}
}
