package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const baseUrl = "http://api.mangadex.org/"
const coverUrl = "https://uploads.mangadex.org/covers/"

func main() {
	router := gin.Default()

	router.GET("/manga", getManga)

	router.Run("localhost:8080")
}

func getManga(c *gin.Context) {
	params := url.Values{}
	params.Add("limit", "10")
	params.Add("includedTagsMode", "AND")
	params.Add("excludedTagsMode", "OR")
	params.Add("availableTranslatedLanguage[]", "en")
	params.Add("contentRating[]", "safe")
	params.Add("contentRating[]", "suggestive")
	params.Add("contentRating[]", "erotica")
	params.Add("order[latestUploadedChapter]", "desc")
	params.Add("includes[]", "manga")
	params.Add("hasAvailableChapters", "true")

	if c.Query("title") != "" {
		params.Add("title", c.Query("title"))
	}

	fullUrl := fmt.Sprintf("%smanga?%s", baseUrl, params.Encode())

	res, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var result MangaResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding JSON for manga:", err)
		return
	}

	var list []Manga
	for _, i := range result.Data {
		var tagList []string
		for _, j := range i.Attributes.Tags {
			tagList = append(tagList, j.Attributes.Name.EN)
		}

		var coverID string
		for _, j := range i.Relationships {
			if j.Type == "cover_art" {
				coverID = j.ID
				break
			}
		}

		fileUrl := fmt.Sprintf("%scover/%s?includes[]=manga", baseUrl, coverID)
		res, err := http.Get(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		var result CoverResponse
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			fmt.Println("Error decoding JSON for cover:", err, fileUrl)
			return
		}

		coverUrl := fmt.Sprintf("%s%s/%s", coverUrl, i.ID, result.Data.Attributes.FileName)

		manga := Manga{
			ID:          i.ID,
			Title:       i.Attributes.Title.EN,
			Description: i.Attributes.Description.EN,
			Created:     i.Attributes.CreatedAt,
			Updated:     i.Attributes.UpdatedAt,
			Status:      i.Attributes.Status,
			Tags:        tagList,
			CoverUrl:    coverUrl,
		}
		list = append(list, manga)
	}

	c.IndentedJSON(http.StatusOK, list)
}
