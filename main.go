package main

import (
	"net/http"

	"fmt"

	"github.com/AhmetTK4/goshort/service"
	"github.com/AhmetTK4/goshort/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	storage.InitRedis()
	r := gin.Default()

	r.POST("/api/shorten", func(c *gin.Context) {
		var request struct {
			URL string `json:"url" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("JSON parse hatası:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		shortCode := service.GenerateShortCode(6)

		err := storage.RDB.Set(storage.Ctx, shortCode, request.URL, 0).Err()
		if err != nil {
			fmt.Println("URL değeri:", request.URL)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_url": "http://localhost:8080/g/" + shortCode,
		})
	})

	r.Run(":8080")
}
