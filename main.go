package main

import (
	"fmt"
	"os"

	"github.com/AvaruusAnkka/manga-backend/src/mangadex"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.SetTrustedProxies(nil)

	router.GET("/manga", mangadex.GetManga)

	router.GET("/chapter", mangadex.GetChapter)

	router.Run(fmt.Sprintf(":%s", port))
}
