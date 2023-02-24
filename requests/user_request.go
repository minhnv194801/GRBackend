package requests

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}
