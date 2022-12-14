package handlers

import (
	"fmt"

	"github.com/gabriel/gabrielyea/go-bank/token"
	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var CurrentServer Server

type Server struct {
	Router     *gin.Engine
	TokenMaker token.Maker
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

	SetRoutes(router, h, tMaker)

	s := &Server{
		Router:     router,
		TokenMaker: tMaker,
		config:     config,
	}

	CurrentServer = *s
	return s
}

func RunServer(config util.Config, h HandlersInt) {
	serv := SetUpServer(config, h)
	serv.Router.Run()
}
