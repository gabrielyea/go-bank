package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gabriel/gabrielyea/go-bank/repo"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (h *handler) CreateAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := repo.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := h.Store.CreateAccount(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account": account,
	})
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (h *handler) GetAccount(c *gin.Context) {
	var req getAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := h.Store.GetAccount(c, req.ID)
	if err != nil {
		errMsg := err
		status := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			errMsg = errors.New(fmt.Sprintf("no account with id: %v found", req.ID))
		}
		c.JSON(status, errorResponse(errMsg))
		return
	}

	c.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (h *handler) ListAccounts(c *gin.Context) {
	var req listAccountRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := repo.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}

	accountList, err := h.Store.ListAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"list": accountList,
	})

}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (h *handler) DeleteAccount(c *gin.Context) {
	var req deleteAccountRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.Store.DeleteAccount(c, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("account with id: %v deleted!", req.ID),
	})
}

type updateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required"`
}

func (h *handler) UpdateAccount(c *gin.Context) {
	var req updateAccountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := repo.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	account, err := h.Store.UpdateAccount(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "account updated!",
		"account": account,
	})
}
