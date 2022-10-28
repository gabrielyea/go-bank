package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gabriel/gabrielyea/go-bank/repo"
	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/gin-gonic/gin"
)

type UserInterface interface {
	CreateUser(*gin.Context)
	LogIn(*gin.Context)
}

type createUserRequest struct {
	UserName string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	UserName string `json:"user_name" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func createUserResponse(user repo.User) userResponse {
	return userResponse{
		UserName: user.UserName,
		FullName: user.FullName,
		Email:    user.Email,
	}
}

func (h *handler) CreateUser(c *gin.Context) {

	var req createUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hash, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := repo.CreateUserParams{
		UserName:       req.UserName,
		HashedPassword: hash,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := h.Store.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := createUserResponse(user)

	c.JSON(http.StatusOK, response)
}

type logInRequest struct {
	UserName string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type logInResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (h *handler) LogIn(c *gin.Context) {
	var req logInRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.Store.GetUser(c, req.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.IsValidPassword(req.Password, user.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := CurrentServer.TokenMaker.CreateToken(
		req.UserName,
		CurrentServer.config.TokenDuration,
	)
	fmt.Printf("accessToken: %v\n", accessToken)

	c.JSON(http.StatusOK, logInResponse{
		AccessToken: accessToken,
		User:        createUserResponse(user),
	})
}
