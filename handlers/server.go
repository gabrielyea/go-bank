package handlers

import (
	"fmt"

	"github.com/gabriel/gabrielyea/go-bank/token"
	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func SetUpServer(config util.Config, h HandlersInt) *Server {
	tMaker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
		return nil
	}
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
		Router:     router,
		tokenMaker: tMaker,
		config:     config,
	}
}

func RunServer(config util.Config, h HandlersInt) {
	serv := SetUpServer(config, h)
	serv.Router.Run()
}
