package paymentcontroller

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"magna/requests"
	"magna/services/paymentservice"
	"magna/services/sessionservice"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type order struct {
	chapterId string
	userId    string
}

var (
	orderMap          map[string]order = make(map[string]order)
	urlMap            map[order]string = make(map[order]string)
	momoExpiredHour                    = 2
	momoExpiredMinute                  = 0
)

func GetMomoPayURLForChapter(c *gin.Context) {
	chapterId := c.Param("chapterid")

	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	req := requests.GetMomoPayURLRequest{}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}

	if urlMap[order{chapterId, userId}] != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"payUrl": urlMap[order{chapterId, userId}]})
		return
	}

	scheme := "https:"
	momoIpnUrl := scheme + "//" + c.Request.Host + "/api/v1/pay/momo/ipn"
	orderId, payUrl, err := paymentservice.GetMomoPayURLForChapter(chapterId, req.RedirectUrl, momoIpnUrl)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}
	newOrder := order{chapterId, userId}
	orderMap[orderId] = newOrder
	urlMap[newOrder] = payUrl
	time.AfterFunc(time.Duration(momoExpiredHour)*time.Hour+time.Duration(momoExpiredMinute)*time.Minute, func() { delete(orderMap, orderId) })
	time.AfterFunc(time.Duration(momoExpiredHour)*time.Hour+time.Duration(momoExpiredMinute)*time.Minute, func() { delete(urlMap, newOrder) })

	c.IndentedJSON(http.StatusOK, gin.H{"payUrl": payUrl})
}

func SetOwned(c *gin.Context) {
	req := requests.MomoIPNRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}
	log.Println("Receive momo ipn request: ", req)

	//build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("partnerCode=")
	rawSignature.WriteString(req.PartnerCode)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(req.OrderId)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(req.RequestId)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(strconv.Itoa(req.Amount))
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(req.OrderInfo)
	rawSignature.WriteString("&orderType=")
	rawSignature.WriteString(req.OrderType)
	rawSignature.WriteString("&transId=")
	rawSignature.WriteString(strconv.Itoa(req.TransId))
	rawSignature.WriteString("&resultCode=")
	rawSignature.WriteString(strconv.Itoa(req.ResultCode))
	rawSignature.WriteString("&message=")
	rawSignature.WriteString(req.Message)
	rawSignature.WriteString("&payType=")
	rawSignature.WriteString(req.PayType)
	rawSignature.WriteString("&responseTime=")
	rawSignature.WriteString(strconv.Itoa(req.ResponseTime))
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(req.ExtraData)

	var publicKey = os.Getenv("MOMO_API_SECRET_KEY")
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(publicKey))

	// Write Data to it
	h.Write(rawSignature.Bytes())
	calculated := h.Sum(nil)
	check := hmac.Equal(calculated, []byte(req.Signature))
	if !check {
		log.Println("ERROR: Fail signature check")
		// c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Fail signature check"})
		// return
	}

	if req.ResultCode != 0 {
		log.Println("ERROR: Result Code: " + strconv.Itoa(req.ResultCode))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Fail payment"})
		return
	}

	paymentservice.SetOwned(orderMap[req.OrderId].userId, orderMap[req.OrderId].chapterId)
	delete(orderMap, req.OrderId)

	c.JSON(http.StatusNoContent, gin.H{})
}
