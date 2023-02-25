package responses

type LoginResponse struct {
	Sessionkey string `json:"sessionkey"`
	Refreshkey string `json:"refreshkey"`
	Id         string `json:"id"`
	IsLogin    bool   `json:"isLogin"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
}

type RegisterResponse struct {
	Sessionkey string `json:"sessionkey"`
	Refreshkey string `json:"refreshkey"`
	Id         string `json:"id"`
	IsLogin    bool   `json:"isLogin"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
}

type RefreshResponse struct {
	Sessionkey string `json:"sessionkey"`
	Refreshkey string `json:"refreshkey"`
	Id         string `json:"id"`
	IsLogin    bool   `json:"isLogin"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
}
