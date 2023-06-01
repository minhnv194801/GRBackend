package responses

type SearchItem struct {
	Id          string          `json:"id"`
	Title       string          `json:"title"`
	Cover       string          `json:"cover"`
	Description string          `json:"description"`
	Status      int             `json:"status"`
	Rating      float32         `json:"rating"`
	Tags        []string        `json:"tags"`
	ChapterList []NewestChapter `json:"chapters"`
}

type SearchResponse struct {
	Data       []SearchItem `json:"data"`
	TotalCount int          `json:"totalCount"`
}
