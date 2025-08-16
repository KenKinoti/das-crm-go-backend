package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetDocuments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get documents endpoint - TODO: implement"})
}

func (h *Handler) GetDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get document endpoint - TODO: implement"})
}

func (h *Handler) UploadDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Upload document endpoint - TODO: implement"})
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update document endpoint - TODO: implement"})
}

func (h *Handler) DeleteDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete document endpoint - TODO: implement"})
}

func (h *Handler) DownloadDocument(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Download document endpoint - TODO: implement"})
}
