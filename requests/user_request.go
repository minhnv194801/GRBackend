package requests

type UpdateUserInfoRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Gender    int    `json:"gender"`
}
