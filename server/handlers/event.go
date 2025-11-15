package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



func (h *Handler) GetEvent(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "GetEvent",
	})
}

func (h *Handler) AddEvent(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "AddEvent",
	})
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "AddEvent",
	})
}