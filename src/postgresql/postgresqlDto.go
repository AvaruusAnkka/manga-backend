package postgresql

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Manga struct {
	ID             string `json:"id"`
	User_id        string `json:"user_id"`
	Manga_id       string `json:"manga_id"`
	LastViewedChap int    `json:"lastViewedChap"`
	Progress       int    `json:"progress"`
}
