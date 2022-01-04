package middlewares

import (
	"errors"
	"fmt"
	"github.com/PengShaw/go-common/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

const (
	UserIDContextKey   = "current_user_id"
	UsernameContextKey = "current_username"
)

var (
	errTokenFormat    = errors.New("token format is wrong")
	errTokenExpired   = errors.New("token expired")
	errTokenNoActive  = errors.New("token not active")
	errTokenUnHandler = errors.New("couldn't handle this token")
)

var (
	privateKey        string
	jwtExpired        int
	jwtTokenHeaderKey string
)

func InitJwt(PrivateKey string, JwtExpired int, HeaderKey *string) {
	privateKey = PrivateKey
	jwtExpired = JwtExpired
	jwtTokenHeaderKey = "JWT-Token"
	if HeaderKey != nil {
		jwtTokenHeaderKey = *HeaderKey
	}
}

func LoginHandler(verifyFunc func(username, password string) (uint, string, bool)) func(*gin.Context) {
	type userLogin struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var jsonSchema userLogin
		if err := c.ShouldBindJSON(&jsonSchema); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		userID, username, ok := verifyFunc(jsonSchema.Username, jsonSchema.Password)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "wrong username or password"})
			return
		}
		logger.Infof("user (%d)[%s] logined", userID, username)

		token, err := getJwt(userID, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func JwtTokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(jwtTokenHeaderKey)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "missing " + jwtTokenHeaderKey + " in Header."})
			return
		}

		claims, err := verifyJwt(tokenString)
		if err != nil {
			if errors.Is(err, errTokenFormat) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
				return
			} else if errors.Is(err, errTokenExpired) || errors.Is(err, errTokenNoActive) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}
		}
		c.Set(UserIDContextKey, claims.ID)
		c.Set(UsernameContextKey, claims.Username)
		c.Next()
	}
}

type jwtClaims struct {
	ID       uint
	Username string
	jwt.StandardClaims
}

func getJwt(id uint, username string) (string, error) {
	key := []byte(privateKey)

	claims := jwtClaims{
		id,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(jwtExpired) * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func verifyJwt(tokenString string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errTokenFormat
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errTokenNoActive
			}
		}
		return nil, fmt.Errorf("%w: [%w]", errTokenUnHandler, err)
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errTokenUnHandler
}
