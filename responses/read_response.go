package responses

import "magna/model"

type ReadChapterListItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type ReadResponse struct {
	MangaId     string          `json:"mangaId"`
	MangaTitle  string          `json:"mangaTitle"`
	Title       string          `json:"title"`
	ChapterList []model.Chapter `json:"chapterList"`
	Pages       []string        `json:"pages"`
}
