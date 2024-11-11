package manga

type Manga struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Chapters    int      `json:"chapters"`
	CoverUrl    string   `json:"coverUrl"`
	Created     string   `json:"created"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
	Updated     string   `json:"updated"`
}

type Chapter struct {
	Id        string   `json:"id"`
	Number    int      `json:"number"`
	ImageUrls []string `json:"imageUrls"`
}
