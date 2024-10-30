package mangadex

type MangaResponse struct {
	Data []Data `json:"data"`
}

type Data struct {
	ID         string `json:"id"`
	Attributes struct {
		Title       en     `json:"title"`
		Description en     `json:"description"`
		Status      string `json:"status"`
		CreatedAt   string `json:"createdat"`
		UpdatedAt   string `json:"updatedat"`
		Tags        []tags `json:"tags"`
	} `json:"attributes"`
	Relationships []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
}

type CoverResponse struct {
	Data struct {
		Attributes struct {
			FileName string `json:"filename"`
		} `json:"attributes"`
	} `json:"data"`
}

type Chapter struct {
	Id        string   `json:"id"`
	Number    int      `json:"number"`
	ImageUrls []string `json:"imageUrls"`
}

type tags struct {
	Attributes struct {
		Name en `json:"name"`
	} `json:"attributes"`
}

type en struct {
	EN string `json:"en"`
}

type Manga struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Chapters    string   `json:"chapters"`
	CoverUrl    string   `json:"coverUrl"`
	Created     string   `json:"created"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
	Updated     string   `json:"updated"`
}
