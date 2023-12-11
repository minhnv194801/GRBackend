package paymentservice

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"magna/model"
	"magna/services/chapterservice"
	"magna/services/userservice"
	"net/http"
	"os"
	"strconv"

	"github.com/sony/sonyflake"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//define a payload, reference in https://developers.momo.vn/#cong-thanh-toan-momo-phuong-thuc-thanh-toan
type Payload struct {
	PartnerCode  string `json:"partnerCode"`
	AccessKey    string `json:"accessKey"`
	RequestID    string `json:"requestId"`
	Amount       string `json:"amount"`
	OrderID      string `json:"orderId"`
	OrderInfo    string `json:"orderInfo"`
	PartnerName  string `json:"partnerName"`
	StoreId      string `json:"storeId"`
	OrderGroupId string `json:"orderGroupId"`
	Lang         string `json:"lang"`
	AutoCapture  bool   `json:"autoCapture"`
	RedirectUrl  string `json:"redirectUrl"`
	IpnUrl       string `json:"ipnUrl"`
	ExtraData    string `json:"extraData"`
	RequestType  string `json:"requestType"`
	Signature    string `json:"signature"`
}

func GetMomoPayURLForChapter(chapterId, redirectUrl, ipnUrl string) (orderId, payUrl string, err error) {
	chapter := new(model.Chapter)
	chapterObjId, err := primitive.ObjectIDFromHex(chapterId)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	err = chapter.GetItemFromObjectId(chapterObjId)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	manga := new(model.Manga)
	err = manga.GetItemFromObjectId(chapter.Manga)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	orderInfo := "Mua chương truyện " + chapter.Name + " của bộ truyện " + manga.Name

	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	//randome orderID and requestID
	a, _ := flake.NextID()
	b, _ := flake.NextID()

	var endpoint = os.Getenv("MOMO_API_URL")
	var accessKey = os.Getenv("MOMO_API_ACCESS_KEY")
	var secretKey = os.Getenv("MOMO_API_SECRET_KEY")
	var partnerCode = os.Getenv("MOMO_API_PARTNER_CODE")
	orderId = strconv.FormatUint(a, 16)
	var requestId = strconv.FormatUint(b, 16)
	var partnerName = ""
	var storeId = ""
	var orderGroupId = ""
	var autoCapture = true
	var lang = "vi"
	var requestType = "payWithMethod"

	//build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(accessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(strconv.Itoa(int(chapter.Price)))
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString("")
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(ipnUrl)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderId)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(partnerCode)
	rawSignature.WriteString("&redirectUrl=")
	rawSignature.WriteString(redirectUrl)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestId)
	rawSignature.WriteString("&requestType=")
	rawSignature.WriteString(requestType)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	hmac := hmac.New(sha256.New, []byte(secretKey))

	// Write Data to it
	hmac.Write(rawSignature.Bytes())
	log.Println("Raw signature: " + rawSignature.String())

	// Get result and encode as hexadecimal string
	signature := hex.EncodeToString(hmac.Sum(nil))

	var payload = Payload{
		PartnerCode:  partnerCode,
		AccessKey:    accessKey,
		RequestID:    requestId,
		Amount:       strconv.Itoa(int(chapter.Price)),
		RequestType:  requestType,
		RedirectUrl:  redirectUrl,
		IpnUrl:       ipnUrl,
		OrderID:      orderId,
		StoreId:      storeId,
		PartnerName:  partnerName,
		OrderGroupId: orderGroupId,
		AutoCapture:  autoCapture,
		Lang:         lang,
		OrderInfo:    orderInfo,
		ExtraData:    "",
		Signature:    signature,
	}

	var jsonPayload []byte
	jsonPayload, err = json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	log.Println("Payload: " + string(jsonPayload))
	log.Println("Signature: " + signature)

	//send HTTP to momo endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalln(err)
		return "", "", err
	}

	//result
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println("Response from Momo: ", result)

	log.Println()
	log.Println()
	log.Println()

	if result["resultCode"].(float64) != 0 {
		return "", "", errors.New("internal server error")
	}

	return orderId, result["payUrl"].(string), nil
}

func SetOwned(userId, chapterId string) error {
	chapter, err := chapterservice.GetChapterInfo(chapterId)
	if err != nil {
		return err
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		return err
	}

	user.OwnedChapters = append(user.OwnedChapters, chapter.Id)
	chapter.OwnedUsers = append(chapter.OwnedUsers, user.Id)

	err = chapter.SetOwned()
	if err != nil {
		return err
	}
	return user.SetOwned()
}
