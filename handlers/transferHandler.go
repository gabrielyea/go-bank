package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gabriel/gabrielyea/go-bank/middleware"
	"github.com/gabriel/gabrielyea/go-bank/repo"
	"github.com/gabriel/gabrielyea/go-bank/token"
	"github.com/gin-gonic/gin"
)

type TransferInt interface {
	CreateTransfer(*gin.Context)
}

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (h *handler) CreateTransfer(c *gin.Context) {
	var req createTransferRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload := c.MustGet(middleware.AuthKeys()["payloadKey"]).(*token.Payload)
	account, ok := h.validAccount(c, req.FromAccountID, req.Currency)
	if !ok {
		return
	}
	if account.Owner != payload.UserName {
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("transfer account doesnt belong to owner")))
		return
	}
	_, ok = h.validAccount(c, req.ToAccountID, req.Currency)
	if !ok {
		return
	}

	arg := repo.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	res, err := h.TransferTx(c, arg)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) validAccount(c *gin.Context, id int64, currency string) (repo.Account, bool) {
	var account repo.Account
	account, err := h.Store.GetAccount(c, id)
	if err != nil {
		errMsg := err
		status := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			errMsg = errors.New(fmt.Sprintf("no account with id: %v found", id))
		}
		c.JSON(status, errorResponse(errMsg))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("currency mismatch, account(%v) selected(%v)", account.Currency, currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, account.Currency == currency
}
