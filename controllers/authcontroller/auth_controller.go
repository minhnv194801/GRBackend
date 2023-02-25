package authcontroller

import (
	"magna/requests"
	"magna/responses"
	"magna/services/authservice"
	"magna/services/sessionservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	req := requests.UserLoginRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}

	sskey, refreshkey, userId, username, avatar, err := authservice.Login(req.Email, req.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, responses.LoginResponse{
		Sessionkey: sskey,
		Refreshkey: refreshkey,
		Id:         userId,
		IsLogin:    true,
		Username:   username,
		Avatar:     avatar,
	})
}

func Register(c *gin.Context) {
	req := requests.UserRegisterRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}

	sskey, refreshkey, userId, username, avatar, err := authservice.Register(req.Email, req.Password, req.RePassword)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, responses.RegisterResponse{
		Sessionkey: sskey,
		Refreshkey: refreshkey,
		Id:         userId,
		IsLogin:    true,
		Username:   username,
		Avatar:     avatar,
	})
}

func RefreshSession(c *gin.Context) {
	refreshkey := c.GetHeader("Authorization")

	sskey, refreshkey, userId, username, avatar, err := sessionservice.RefreshSession(refreshkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, responses.RefreshResponse{
		Sessionkey: sskey,
		Refreshkey: refreshkey,
		Id:         userId,
		IsLogin:    true,
		Username:   username,
		Avatar:     avatar,
	})
}
