package main

import (
	"fmt"
	"os"

	"github.com/AvaruusAnkka/manga-backend/src/mangadex"
	"github.com/AvaruusAnkka/manga-backend/src/postgresql"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/manga", mangadex.GetManga)

	router.GET("/chapter", mangadex.GetChapter)

	router.GET("/user/id", postgresql.GetUser)
	router.GET("/user/add", postgresql.AddUser)

	router.Run(fmt.Sprintf(":%s", port))
}
