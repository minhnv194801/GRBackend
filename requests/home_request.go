package requests

type HotItemsListRequest struct {
	Count int `json:"count"`
}

type RecommendListRequest struct {
	Count int `json:"count"`
}

type NewestListRequest struct {
	Postition int `json:"position"`
	Count     int `json:"count"`
}
