package handlers

import (
	"github.com/gabriel/gabrielyea/go-bank/repo"
	"github.com/gin-gonic/gin"
)

type HandlersInt interface {
	AccountInt
	TransferInt
	UserInterface
}

type handler struct {
	repo.Store
}

func NewHandler(r repo.Store) HandlersInt {
	return &handler{r}
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
