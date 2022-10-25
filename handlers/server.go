package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Router *gin.Engine
}

func SetUpServer(h HandlersInt) *Server {
	server := gin.Default()
	router := server

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", h.CreateAccount)
	router.POST("/transfers", h.CreateTransfer)
	router.POST("/users", h.CreateUser)
	router.GET("/accounts/:id", h.GetAccount)
	router.GET("accounts", h.ListAccounts)
	router.DELETE("accounts/:id", h.DeleteAccount)
	router.PATCH("/accounts", h.UpdateAccount)

	return &Server{
		Router: router,
	}
}

func RunServer(h HandlersInt) {
	serv := SetUpServer(h)
	serv.Router.Run()
}
