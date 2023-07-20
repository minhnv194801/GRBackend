package responses

type UserReferenceItem struct {
	DisplayName string `json:"displayname"`
	Avatar      string `json:"avatar"`
}

type MangaReferenceItem struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
}

type ChapterReferenceItem struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
}

type CommentReferenceItem struct {
	Content     string             `json:"content"`
	TimeCreated uint               `json:"timeCreated"`
	User        UserReferenceItem  `json:"user"`
	Manga       MangaReferenceItem `json:"manga"`
}

type ReportReferenceItem struct {
	Content     string               `json:"content"`
	TimeCreated uint                 `json:"timeCreated"`
	User        UserReferenceItem    `json:"user"`
	Chapter     ChapterReferenceItem `json:"chapter"`
	Status      int                  `json:"status"`
}
