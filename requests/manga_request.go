package requests

type MangaChapterListRequest struct {
	Postition int `json:"position"`
	Count     int `json:"count"`
}

type CommentListRequest struct {
	Postition int `json:"position"`
	Count     int `json:"count"`
}
