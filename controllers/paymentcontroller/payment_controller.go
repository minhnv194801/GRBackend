package controllers

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

	"github.com/gin-gonic/gin"
)

func GetMomoPayURL(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	_ = userId

	req := requests.GetMomoPayURLRequest{}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}

	// TODO: link orderId with userId+chapterId
	orderId, payUrl, err := paymentservice.GetMomoPayURL(req.OrderInfo, req.RedirectUrl, "ipnUrl", strconv.Itoa(req.Amount), req.ExtraData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}
	_ = orderId

	c.IndentedJSON(http.StatusOK, gin.H{"payUrl": payUrl})
}

// http://localhost:8081/?partnerCode=MOMO&orderId=6cd25bf88002596&requestId=6cd25bf88012596&amount=10000&orderInfo=pay+with+MoMo&orderType=momo_wallet&transId=3105694606&resultCode=0&message=Successful.&payType=napas&responseTime=1701645754029&extraData=&signature=7c7f9acccaa791344cf63a9819ca594871b05d51bbff4f798d6e72beaa3c59f2
// TODO: check hash
func SetOwned(c *gin.Context) {
	partnerCode := c.Param("partnerCode")
	orderId := c.Param("orderId")
	requestId := c.Param("requestId")
	amount := c.Param("amount")
	orderInfo := c.Param("orderInfo")
	orderType := c.Param("orderType")
	transId := c.Param("transId")
	resultCode := c.Param("resultCode")
	message := c.Param("message")
	payType := c.Param("payType")
	responseTime := c.Param("responseTime")
	extraData := c.Param("extraData")
	signature := c.Param("signature")

	//build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("partnerCode=")
	rawSignature.WriteString(partnerCode)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amount)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&orderType=")
	rawSignature.WriteString(orderType)
	rawSignature.WriteString("&transId=")
	rawSignature.WriteString(transId)
	rawSignature.WriteString("&resultCode=")
	rawSignature.WriteString(resultCode)
	rawSignature.WriteString("&message=")
	rawSignature.WriteString(message)
	rawSignature.WriteString("&payType=")
	rawSignature.WriteString(payType)
	rawSignature.WriteString("&responseTime=")
	rawSignature.WriteString(responseTime)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)

	var publicKey = os.Getenv("MOMO_API_PUBLIC_KEY")
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(publicKey))

	// Write Data to it
	h.Write(rawSignature.Bytes())
	calculated := h.Sum(nil)
	check := hmac.Equal(calculated, []byte(signature))
	if !check {
		log.Println("ERROR: Fail signature check")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Fail signature check"})
		return
	}

	if resultCode != "0" {
		log.Println("ERROR: Result Code: " + resultCode)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Fail signature check"})
		return
	}

	//TODO: get userId and chapterId from orderId

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}
