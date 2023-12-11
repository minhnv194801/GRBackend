package requests

type GetMomoPayURLRequest struct {
	RedirectUrl string `json:"redirectUrl"`
}

type MomoIPNRequest struct {
	PartnerCode  string `json:"partnerCode"`
	OrderId      string `json:"orderId"`
	RequestId    string `json:"requestId"`
	Amount       int    `json:"amount"`
	OrderInfo    string `json:"orderInfo"`
	OrderType    string `json:"orderType"`
	TransId      int    `json:"transId"`
	ResultCode   int    `json:"resultCode"`
	Message      string `json:"message"`
	PayType      string `json:"payType"`
	ResponseTime int    `json:"responseTime"`
	ExtraData    string `json:"extraData"`
	Signature    string `json:"signature"`
}
