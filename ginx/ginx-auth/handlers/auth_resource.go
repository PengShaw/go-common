package handlers

import (
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListAuthResource(c *gin.Context) {
	items, err := models.ListAuthResource()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	type resp struct {
		ID     uint   `json:"id"`
		Uri    string `json:"uri"`
		Method string `json:"method"`
	}
	var res []resp
	for _, v := range items {
		item := resp{ID: v.ID, Uri: v.Uri, Method: v.Method}
		res = append(res, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
