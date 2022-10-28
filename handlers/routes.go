package handlers

import (
	"github.com/gabriel/gabrielyea/go-bank/middleware"
	"github.com/gabriel/gabrielyea/go-bank/token"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine, h HandlersInt, tm token.Maker) {
	v1 := r.Group("/")
	{
		publicR := v1.Group("/")
		{
			publicRoutes(publicR, h)
		}
		privateR := v1.Group("/")
		{
			privateRoutes(privateR, h, tm)
		}
	}
}

func publicRoutes(r *gin.RouterGroup, h HandlersInt) {
	{
		r.POST("/users", h.CreateUser)
		r.POST("/login", h.LogIn)
	}
}

func privateRoutes(r *gin.RouterGroup, h HandlersInt, tm token.Maker) {

	r.Use(middleware.AuthMiddleware(tm))
	{
		r.POST("/accounts", h.CreateAccount)
		r.GET("/accounts/:id", h.GetAccount)
		r.GET("/accounts", h.ListAccounts)
		r.POST("/transfers", h.CreateTransfer)
		r.PATCH("/accounts", h.UpdateAccount)
		r.DELETE("/accounts/:id", h.DeleteAccount)
	}
}
