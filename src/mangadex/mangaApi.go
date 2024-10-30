package mangadex

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
)

const baseUrl = "http://api.mangadex.org/"
const coverUrl = "https://uploads.mangadex.org/covers/"

func GetManga(c *gin.Context) {
	fullUrl := baseUrl + "manga?" + getParams(c)

	res, err := http.Get(fullUrl)
	if err != nil {
		c.String(400, err.Error())
	}
	defer res.Body.Close()

	var mangaResponse MangaResponse
	if err := json.NewDecoder(res.Body).Decode(&mangaResponse); err != nil {
		c.String(400, "Error decoding JSON for manga:"+err.Error())
	}

	list := mangaConverter(mangaResponse)
	getCover(list)

	c.IndentedJSON(http.StatusOK, list)
}

func getParams(c *gin.Context) string {
	params := url.Values{}
	params.Add("limit", "10")
	params.Add("availableTranslatedLanguage[]", "en")
	params.Add("order[latestUploadedChapter]", "desc")
	params.Add("hasAvailableChapters", "true")
	params.Add("title", c.DefaultQuery("title", ""))

	return params.Encode()
}

func getTags(data *Data) []string {
	var tagList []string
	for _, tag := range data.Attributes.Tags {
		tagList = append(tagList, tag.Attributes.Name.EN)
	}

	return tagList
}

func getCoverUrl(data *Data) string {
	var coverId string
	for _, relationship := range data.Relationships {
		if relationship.Type == "cover_art" {
			coverId = relationship.ID
			break
		}
	}

	return baseUrl + "cover/" + coverId + "?includes[]=manga"
}

func mangaConverter(res MangaResponse) []Manga {
	var list []Manga
	for i, data := range res.Data {
		list = append(list, Manga{
			ID:          data.ID,
			Title:       data.Attributes.Title.EN,
			Description: data.Attributes.Description.EN,
			Created:     data.Attributes.CreatedAt,
			Updated:     data.Attributes.UpdatedAt,
			Status:      data.Attributes.Status,
			Tags:        getTags(&res.Data[i]),
			CoverUrl:    getCoverUrl(&res.Data[i]),
		})
	}

	return list
}

func getCover(mangas []Manga) {
	var wg sync.WaitGroup
	for i, manga := range mangas {
		wg.Add(1)

		go func(manga Manga) {
			defer wg.Done()

			res, err := http.Get(manga.CoverUrl)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer res.Body.Close()

			var result CoverResponse
			if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
				fmt.Println("Error decoding JSON for cover:", err, manga.CoverUrl)
				return
			}

			mangas[i].CoverUrl = coverUrl + manga.ID + "/" + result.Data.Attributes.FileName
		}(manga)

		wg.Wait()
	}
}
