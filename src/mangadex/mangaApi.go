package mangadex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/AvaruusAnkka/manga-backend/src/manga"
	"github.com/gin-gonic/gin"
)

const baseUrl = "http://api.mangadex.org/"
const coverUrl = "https://uploads.mangadex.org/"

func GetManga(c *gin.Context) {
	url := baseUrl + "manga?" + getParams(c)

	res, err := http.Get(url)
	if err != nil {
		c.String(400, err.Error())
	}
	defer res.Body.Close()

	var mangaResponse MangaResponse
	if err := json.NewDecoder(res.Body).Decode(&mangaResponse); err != nil {
		c.String(400, err.Error())
	}

	list := mangaConverter(mangaResponse)

	c.IndentedJSON(http.StatusOK, list)
}

func GetChapter(c *gin.Context) {
	id := c.Query("id")
	chapterNumber := c.Query("chapter")

	number, err := strconv.Atoi(chapterNumber)

	if id == "" || chapterNumber == "" || err != nil {
		c.String(400, "Missing/Invalid queries")
		return
	}

	chapterId, err := getChapterId(id, chapterNumber)
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}

	chapter := manga.Chapter{
		Id:        chapterId,
		Number:    number,
		ImageUrls: getChapterUrls(chapterId),
	}

	c.IndentedJSON(200, chapter)
}

func getChapterUrls(id string) (list []string) {
	res, _ := http.Get(baseUrl + "at-home/server/" + id)

	var home HomeResponse
	json.NewDecoder(res.Body).Decode(&home)

	for _, file := range home.Chapter.Data {
		fileUrl := coverUrl + "data/" + home.Chapter.Hash + "/" + file
		list = append(list, fileUrl)
	}

	return
}

func getChapterId(id string, chapter string) (chapterId string, err error) {
	url := baseUrl + "manga/" + id + "/aggregate?translatedLanguage%5B%5D=en"
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var aggregate AggregateResponse
	if err = json.NewDecoder(res.Body).Decode(&aggregate); err != nil {
		return
	}

	for _, volume := range aggregate.Volumes {
		for _, number := range volume.Chapters {
			if chapter == number.Chapter {
				chapterId = number.ID
			}
		}
	}

	return
}

func getFileUrls(data HomeResponse) (list []string) {
	for _, file := range data.Chapter.Data {
		fileUrl := coverUrl + "data/" + data.Chapter.Hash + "/" + file
		list = append(list, fileUrl)
	}

	return
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
		fmt.Println(err.Error())
		return ""
	}

	return coverUrl + "covers/" + data.ID + "/" + result.Data.Attributes.FileName
}

func mangaConverter(res MangaResponse) (list []manga.Manga) {
	covers, chapters := getCoverAndChapter(res)

	for i, data := range res.Data {
		list = append(list, manga.Manga{
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
		for _, chapter := range volume.Chapters {
			num, err := strconv.Atoi(chapter.Chapter)
			if err == nil && latestChapter < num {
				latestChapter = num
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
