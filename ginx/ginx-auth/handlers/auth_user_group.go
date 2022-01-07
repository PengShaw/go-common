package handlers

import (
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	contextkey "github.com/PengShaw/go-common/ginx/ginx-context-key"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateUserGroup(c *gin.Context) {
	type itemCreate struct {
		Name    string `json:"name" binding:"required"`
		Comment string `json:"Comment"`
	}

	var jsonSchema itemCreate
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	err := models.CreateUserGroup(jsonSchema.Name, jsonSchema.Comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "create user group successful"})

}

func ListUserGroup(c *gin.Context) {
	groups, err := models.ListUserGroup()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	type resp struct {
		ID        uint   `json:"id"`
		GroupName string `json:"name"`
		Comment   string `json:"comment"`
	}
	var res []resp
	for _, v := range groups {
		item := resp{ID: v.ID, GroupName: v.GroupName, Comment: v.Comment}
		res = append(res, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func DeleteUserGroup(c *gin.Context) {
	type itemDelete struct {
		ID string `uri:"id" binding:"required"`
	}

	var uriSchema = &itemDelete{}
	if err := c.ShouldBindUri(uriSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	id, err := strconv.Atoi(uriSchema.ID)
	if err := c.ShouldBindUri(uriSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	group, err := models.GetUserGroup(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if group.GroupName == contextkey.DefaultAuthAdminUsergroup {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "cannot delete default admin group user"})
		return
	}

	err = models.DeleteUserGroup(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "delete user group successful"})
}

func RenameUserGroup(c *gin.Context) {
	type userUpdate struct {
		UserGroupID uint   `json:"user_group_id" binding:"required"`
		GroupName   string `json:"group_name" binding:"required"`
	}
	var jsonSchema userUpdate
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	group, err := models.GetUserGroup(jsonSchema.UserGroupID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if group.GroupName == contextkey.DefaultAuthAdminUsergroup {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "cannot rename default admin user group"})
		return
	}

	err = models.RenameUserGroup(jsonSchema.UserGroupID, jsonSchema.GroupName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "rename user group successful"})
}

func CommentUserGroup(c *gin.Context) {
	type userUpdate struct {
		UserGroupID uint   `json:"user_group_id" binding:"required"`
		Comment     string `json:"comment" binding:"required"`
	}
	var jsonSchema userUpdate
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	group, err := models.GetUserGroup(jsonSchema.UserGroupID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if group.GroupName == contextkey.DefaultAuthAdminUsergroup {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "cannot comment default admin user group"})
		return
	}

	err = models.ChangeUserGroupComment(jsonSchema.UserGroupID, jsonSchema.Comment)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "comment user group successful"})
}
