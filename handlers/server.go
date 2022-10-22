package handlers

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func SetUpServer(h HandlersInt) *Server {
	server := gin.Default()
	router := server

	router.POST("/accounts", h.CreateAccount)
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
