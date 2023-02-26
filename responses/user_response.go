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
	Name       string `json:"name"`
	UpdateTime uint   `json:"updateTime"`
}

type FavoriteItem struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Cover       string            `json:"cover"`
	ChapterList []FavoriteChapter `json:"chapters"`
}
