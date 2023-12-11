package responses

type ReadChapterListItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type ReadResponse struct {
	MangaId    string   `json:"mangaId"`
	MangaTitle string   `json:"mangaTitle"`
	Title      string   `json:"title"`
	Price      uint     `json:"price"`
	Pages      []string `json:"pages"`
}
