package main

import (
	"net/http"
	"time"

	"github.com/AhmetTK4/goshort/service"
	"github.com/AhmetTK4/goshort/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	storage.InitRedis()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	r.POST("/api/shorten", func(c *gin.Context) {
		var request struct {
			URL string `json:"url" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		shortCodeKey := "original:" + request.URL

		existingShortCode, err := storage.RDB.Get(storage.Ctx, shortCodeKey).Result()
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"short_url": "http://localhost:8080/g/" + existingShortCode,
			})
			return
		}

		shortCode := service.GenerateShortCode(6)

		err = storage.RDB.Set(storage.Ctx, shortCode, request.URL, 1*time.Hour).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		err = storage.RDB.Set(storage.Ctx, shortCodeKey, shortCode, 1*time.Hour).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_url": "http://localhost:8080/g/" + shortCode,
		})
	})

	r.GET("/g/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		longURL, err := storage.RDB.Get(storage.Ctx, shortCode).Result()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}

		storage.RDB.Incr(storage.Ctx, "clicks:"+shortCode)
		c.Redirect(http.StatusFound, longURL)
	})

	r.GET("/api/stats/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		count, err := storage.RDB.Get(storage.Ctx, "clicks:"+shortCode).Result()
		if err != nil {
			count = "0"
		}

		c.JSON(http.StatusOK, gin.H{
			"short_code": shortCode,
			"clicks":     count,
		})
	})

	r.Run(":8080")
}
