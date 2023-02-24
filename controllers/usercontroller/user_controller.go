package usercontroller

import (
	"fmt"
	"magna/requests"
	"magna/responses"
	"magna/services/sessionservice"
	"magna/services/userservice"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	req := requests.UserLoginRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}

	sskey, refreshkey, userId, username, avatar, err := userservice.Login(req.Email, req.Password)
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

	sskey, refreshkey, userId, username, avatar, err := userservice.Register(req.Email, req.Password, req.RePassword)
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

func Test(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	fmt.Println(c.GetHeader("Authorization"))
	token, err := sessionservice.CheckSession(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["id"])
	fmt.Println(id)
	c.IndentedJSON(http.StatusOK, "")
}
