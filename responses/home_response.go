package responses

type NewestChapter struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	UpdateTime uint   `json:"updateTime"`
}

type NewestResponse struct {
	Id          string          `json:"id"`
	Title       string          `json:"title"`
	Cover       string          `json:"cover"`
	ChapterList []NewestChapter `json:"chapters"`
}

type HotItemsResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type RecommendResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}
