package requests

type AdminCreateAccountRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Role     string `json:"Role"`
}

type AdminCreateMangaRequest struct {
	Name          string   `json:"Name"`
	AlternateName []string `json:"AlternateName"`
	Author        string   `json:"Author"`
	Cover         string   `json:"Cover"`
	Description   string   `json:"Description"`
	IsRecommend   bool     `json:"IsRecommend"`
	Tags          []string `json:"Tags"`
}

type AdminCreateChapterRequest struct {
	MangaId string   `json:"Manga"`
	Title   string   `json:"Title"`
	Cover   string   `json:"Cover"`
	Price   uint     `json:"Price"`
	Images  []string `json:"Images"`
}
