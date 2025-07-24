package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.POST("/api/shorten", func(c *gin.Context) {
		var request struct {
			URL string `json:"url" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		shortCode := "abc123"

		c.JSON(http.StatusOK, gin.H{
			"short_url": "http://localhost:8080/g/" + shortCode,
		})
	})

	r.Run(":8080")
}
