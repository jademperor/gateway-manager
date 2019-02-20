package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LocalManageAPIS ...
func LocalManageAPIS(c *gin.Context) {
	c.String(http.StatusOK, "/v1/local/apis")
}
