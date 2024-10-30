package mangadex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const baseUrl = "http://api.mangadex.org/"

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

func getCover(data *Data) string {
	const coverUrl = "https://uploads.mangadex.org/covers/"

	var coverId string
	for _, relationship := range data.Relationships {
		if relationship.Type == "cover_art" {
			coverId = relationship.ID
			break
		}
	}

	fileUrl := baseUrl + "cover/" + coverId + "?includes[]=manga"
	res, err := http.Get(fileUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	var cover CoverResponse
	if err := json.NewDecoder(res.Body).Decode(&cover); err != nil {
		fmt.Println("Error decoding JSON for cover:", err, fileUrl)
	}

	return coverUrl + data.ID + "/" + cover.Data.Attributes.FileName
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
			CoverUrl:    getCover(&res.Data[i]),
		})
	}

	return list
}
