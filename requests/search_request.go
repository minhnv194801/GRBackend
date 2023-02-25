package requests

type SearchRequest struct {
	Query    string   `json:"query"`
	Tags     []string `json:"tags"`
	Position int      `json:"position"`
	Count    int      `json:"count"`
}
