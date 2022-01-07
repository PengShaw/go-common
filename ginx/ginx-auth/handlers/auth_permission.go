package handlers

import (
	"github.com/PengShaw/go-common/ginx/ginx-auth/casbinx"
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	contextkey "github.com/PengShaw/go-common/ginx/ginx-context-key"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateAuthPermission(c *gin.Context) {
	type itemCreate struct {
		GroupID  uint `json:"user_group_id" binding:"required"`
		SourceID uint `json:"resource_id" binding:"required"`
	}

	var jsonSchema itemCreate
	if err := c.ShouldBindJSON(&jsonSchema); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	err := casbinx.AddPermission(jsonSchema.GroupID, jsonSchema.SourceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "add permission successful"})
}

func ListAuthPermission(c *gin.Context) {
	items, err := models.ListPermission()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	type resp struct {
		ID              uint   `json:"id"`
		AuthUserGroupID uint   `json:"group_id"`
		GroupName       string `json:"group_name"`
		AuthResourceID  uint   `json:"resource_id"`
		Uri             string `json:"uri"`
		Method          string `json:"method"`
	}
	var res []resp
	for _, v := range items {
		item := resp{
			ID:              v.ID,
			AuthUserGroupID: v.AuthUserGroupID,
			GroupName:       v.AuthUserGroup.GroupName,
			AuthResourceID:  v.AuthResourceID,
			Uri:             v.AuthResource.Uri,
			Method:          v.AuthResource.Method,
		}
		res = append(res, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func DeleteAuthPermission(c *gin.Context) {
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

	permission, err := models.GetPermissionByID(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if permission.AuthUserGroup.GroupName == contextkey.DefaultAuthAdminUsergroup {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "cannot delete default admin group permission"})
		return
	}

	err = casbinx.DeletePermission(permission.AuthUserGroupID, permission.AuthResourceID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "delete permission successful"})
}
