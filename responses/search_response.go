package responses

type SearchItem struct {
	Id          string          `json:"id"`
	Title       string          `json:"title"`
	Cover       string          `json:"cover"`
	ChapterList []NewestChapter `json:"chapters"`
}

type SearchResponse struct {
	Data       []SearchItem `json:"data"`
	TotalCount int          `json:"totalCount"`
}
