package requests

type GetMomoPayURLRequest struct {
	ChapterId   string `json:"chapterId"`
	OrderInfo   string `json:"orderInfo"`
	RedirectUrl string `json:"redirectUrl"`
	Amount      int    `json:"amount"`
	ExtraData   string `json:"extraData"`
}
