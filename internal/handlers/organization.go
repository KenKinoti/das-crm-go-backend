package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOrganization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get organization endpoint - TODO: implement"})
}

func (h *Handler) UpdateOrganization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update organization endpoint - TODO: implement"})
}
