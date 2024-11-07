package mangadex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
		c.String(400, "Error decoding JSON for manga: "+err.Error())
	}

	covers, chapters := getCoverAndChapter(mangaResponse)
	list := mangaConverter(mangaResponse, covers, chapters)

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

func getTags(data *Data) (tagList []string) {
	for _, tag := range data.Attributes.Tags {
		tagList = append(tagList, tag.Attributes.Name.EN)
	}

	return
}

func getCoverUrl(data *Data) string {
	var id string
	for _, relationship := range data.Relationships {
		if relationship.Type == "cover_art" {
			id = relationship.ID
			break
		}
	}

	url := baseUrl + "cover/" + id + "?includes[]=manga"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer res.Body.Close()

	var result CoverResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding JSON for cover:", err.Error())
		return ""
	}

	return coverUrl + data.ID + "/" + result.Data.Attributes.FileName
}

func mangaConverter(res MangaResponse, covers []string, chapters []int) (list []Manga) {
	for i, data := range res.Data {
		list = append(list, Manga{
			ID:          data.ID,
			Title:       data.Attributes.Title.EN,
			Description: data.Attributes.Description.EN,
			Created:     data.Attributes.CreatedAt,
			Updated:     data.Attributes.UpdatedAt,
			Status:      data.Attributes.Status,
			Tags:        getTags(&res.Data[i]),
			CoverUrl:    covers[i],
			Chapters:    chapters[i],
		})
	}
	return
}

func prepRequests(data *Data, wg *sync.WaitGroup) (cover string, chapter int) {
	defer wg.Done()

	innerWg := sync.WaitGroup{}
	innerWg.Add(2)

	go func() {
		defer innerWg.Done()
		cover = getCoverUrl(data)
	}()

	go func() {
		defer innerWg.Done()
		chapter = getChapterTotal(data.ID)
	}()

	innerWg.Wait()

	return
}

func getChapterTotal(id string) (latestChapter int) {
	url := baseUrl + "manga/" + id + "/aggregate"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var response AggregateResponse
	json.Unmarshal(body, &response)

	for _, volume := range response.Volumes {
		for i := range volume.Chapters {
			number, err := strconv.Atoi(i)
			if err == nil && latestChapter < number {
				latestChapter = number
			}
		}
	}

	return
}

func getCoverAndChapter(res MangaResponse) ([]string, []int) {
	var wg sync.WaitGroup

	covers := make([]string, len(res.Data))
	chapters := make([]int, len(res.Data))

	for i := range res.Data {
		wg.Add(1)

		go func() {
			covers[i], chapters[i] = prepRequests(&res.Data[i], &wg)
		}()
	}

	wg.Wait()

	return covers, chapters
}
