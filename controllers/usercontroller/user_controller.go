package usercontroller

import (
	"bytes"
	"io/ioutil"
	"log"
	"magna/requests"
	"magna/responses"
	"magna/services/sessionservice"
	"magna/services/userservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	var response responses.UserInfoResponse
	response.Email = user.Email
	response.FirstName = user.FirstName
	response.LastName = user.LastName
	response.Gender = user.Gender
	response.Role = user.Role
	c.IndentedJSON(http.StatusOK, response)
}

func UpdateUserInfo(c *gin.Context) {
	var request requests.UpdateUserInfoRequest
	body, _ := ioutil.ReadAll(c.Request.Body)
	println(string(body))

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.DisplayName = request.Username
	user.Avatar = request.Avatar
	user.Gender = request.Gender
	log.Println("request", request.Gender)
	if user.UpdateInfo() != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}
