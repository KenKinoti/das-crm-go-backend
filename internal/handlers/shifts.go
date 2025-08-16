package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetShifts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get shifts endpoint - TODO: implement"})
}

func (h *Handler) GetShift(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get shift endpoint - TODO: implement"})
}

func (h *Handler) CreateShift(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create shift endpoint - TODO: implement"})
}

func (h *Handler) UpdateShift(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update shift endpoint - TODO: implement"})
}

func (h *Handler) UpdateShiftStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update shift status endpoint - TODO: implement"})
}

func (h *Handler) DeleteShift(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete shift endpoint - TODO: implement"})
}
