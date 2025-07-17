package controller

import (
	"boardwallfloor/ckd/internal/db"
	"boardwallfloor/ckd/internal/middleware"
	"boardwallfloor/ckd/internal/service"
	"boardwallfloor/ckd/internal/util"
	"encoding/json"
	"net/http"
)

type TransactionController struct {
	txService service.TransactionService
}

func NewTransactionController(s service.TransactionService) *TransactionController {
	return &TransactionController{txService: s}
}

type ProcessTransactionRequest struct {
	ID     string `json:"id,omitempty"`
	Amount string `json:"amount"`
}

func (c *TransactionController) GetTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.AuthClaimsKey).(*service.AuthClaims)
	if !ok {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	transactions, err := c.txService.ListTransactionsByUserID(r.Context(), claims.UserID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusOK, util.Response{
		Success: true,
		Message: "Transactions fetched successfully",
		Data:    transactions,
	})
}

func (c *TransactionController) ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.AuthClaimsKey).(*service.AuthClaims)
	if !ok {
		util.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	var req ProcessTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.ID != "" {
		util.RespondWithError(w, http.StatusNotImplemented, "Update transaction is not implemented yet")
		return
	}

	params := db.CreateTransactionParams{
		UserID: claims.UserID,
		Amount: req.Amount,
	}

	tx, err := c.txService.CreateTransaction(r.Context(), params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, util.Response{
		Success: true,
		Message: "Transaction created successfully",
		Data:    tx,
	})
}
