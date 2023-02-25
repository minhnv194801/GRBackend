package responses

type UserInfoResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Gender    int    `json:"gender"`
	Role      string `json:"role"`
}
