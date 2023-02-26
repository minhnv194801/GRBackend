package responses

type UserInfoResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Gender    int    `json:"gender"`
	Role      string `json:"role"`
}

type OwnedChapterItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type OwnedChapterResponse struct {
	Id       string             `json:"id"`
	Cover    string             `json:"cover"`
	Title    string             `json:"title"`
	Chapters []OwnedChapterItem `json:"chapters"`
}

type FavoriteChapter struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	UpdateTime uint   `json:"updateTime"`
}

type FavoriteItem struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Cover       string            `json:"cover"`
	ChapterList []FavoriteChapter `json:"chapters"`
}

type ReportResponse struct {
	ChapterId    string `json:"chapterId"`
	ChapterCover string `json:"chapterCover"`
	ChapterTitle string `json:"chapterTitle"`
	TimeCreated  int    `json:"timeCreated"`
	Content      string `json:"content"`
	Status       int    `json:"status"`
	Response     string `json:"response"`
}
