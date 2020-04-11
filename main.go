package main

import (
	"fmt"
	"offerlist/models"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Colly(c *gin.Context) {
	platform_id := c.Query("platform_id")
	short_id := c.Query("short_id")
	fmt.Println(short_id, platform_id)
	c.String(200, models.FetchPrice(platform_id, short_id))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/colly", Colly)
	r.Run(":" + port)
}
