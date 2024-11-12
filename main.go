package main

import (
	"github.com/AvaruusAnkka/manga-backend/src/mangadex"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/manga", mangadex.GetManga)

	router.GET("/chapter", mangadex.GetChapter)

	router.Run("localhost:8080")
}
