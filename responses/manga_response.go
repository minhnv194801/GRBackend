package responses

type MangaInfoResponse struct {
	Title        string   `json:"title"`
	Cover        string   `json:"cover"`
	IsFavorite   bool     `json:"isFavorite"`
	Author       string   `json:"author"`
	Status       uint     `json:"status"`
	Tags         []string `json:"tags"`
	UserRating   uint     `json:"userRating"`
	AvgRating    uint     `json:"avgRating"`
	RatingCount  uint     `json:"ratingCount"`
	Description  string   `json:"description"`
	ChapterCount int      `json:"chapterCount"`
}

type ChapterListResponse struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Cover      string `json:"cover"`
	Price      uint   `json:"price"`
	IsOwned    bool   `json:"isOwned"`
	UpdateTime uint   `json:"updateTime"`
}

type CommentListResponse struct {
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Content    string `json:"content"`
	UpdateTime uint   `json:"updateTime"`
}
