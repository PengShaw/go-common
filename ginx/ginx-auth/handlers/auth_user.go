package handlers

import (
	jwttoken "github.com/PengShaw/go-common/ginx/ginx-auth/jwt-token"
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	contextkey "github.com/PengShaw/go-common/ginx/ginx-context-key"
	"github.com/PengShaw/go-common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LoginHandler(c *gin.Context) {
	type userLogin struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var jsonSchema userLogin
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// query & verify user
	user, err := models.VerifyUser(jsonSchema.Username, jsonSchema.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		return
	}
	logger.Infof("user (%d)[%s] logined", user.ID, user.Username)

	// generate token
	token, err := jwttoken.Generate(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CreateUser(c *gin.Context) {
	type userCreate struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		UserGroupID uint   `json:"user_group_id" binding:"required"`
	}

	var jsonSchema userCreate
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	err := models.CreateUser(jsonSchema.Username, jsonSchema.Password, jsonSchema.UserGroupID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "create user successful"})
}

func ListUser(c *gin.Context) {
	users, err := models.ListUser()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	type resp struct {
		ID          uint   `json:"id"`
		Username    string `json:"username"`
		UserGroupID uint   `json:"user_group_id"`
		UserGroup   string `json:"user_group"`
	}
	var res []resp
	for _, v := range users {
		item := resp{ID: v.ID, Username: v.Username, UserGroupID: v.AuthUserGroupID, UserGroup: v.AuthUserGroup.GroupName}
		res = append(res, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func DeleteUser(c *gin.Context) {
	type userDelete struct {
		ID string `uri:"id" binding:"required"`
	}

	var uriSchema = &userDelete{}
	if err := c.ShouldBindUri(uriSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	id, err := strconv.Atoi(uriSchema.ID)
	if err := c.ShouldBindUri(uriSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	user, err := models.GetUser(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if user.Username == contextkey.DefaultAuthAdminUsername {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "cannot delete default admin user"})
		return
	}

	err = models.DeleteUser(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "delete user successful"})
}

func UserChangePassword(c *gin.Context) {
	type userUpdate struct {
		Password string `json:"password" binding:"required"`
	}
	var jsonSchema userUpdate
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	var userID uint
	userIDValue, ok := c.Get(contextkey.UserIDContextKey)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "missing userID context"})
		return
	}
	if userID, ok = userIDValue.(uint); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "userID context type wrong"})
		return
	}

	err := models.UserChangePassword(userID, jsonSchema.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "change password successful"})
}

func RegroupUser(c *gin.Context) {
	type userUpdate struct {
		UserID      uint `json:"user_id" binding:"required"`
		UserGroupID uint `json:"user_group_id" binding:"required"`
	}
	var jsonSchema userUpdate
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	user, err := models.GetUser(jsonSchema.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if user.Username == contextkey.DefaultAuthAdminUsername {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "cannot regroup default admin user"})
		return
	}

	err = models.RegroupUser(jsonSchema.UserID, jsonSchema.UserGroupID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "regroup user successful"})
}
