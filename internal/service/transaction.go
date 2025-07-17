package service

import (
	"boardwallfloor/ckd/internal/db"
	utils "boardwallfloor/ckd/internal/util"
	"context"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, params db.CreateTransactionParams) (db.Transaction, error)
	ListTransactionsByUserID(ctx context.Context, userID int32) ([]db.Transaction, error)
}

type transactionService struct {
	queries *db.Queries
}

func NewTransactionService(queries *db.Queries) TransactionService {
	return &transactionService{queries: queries}
}

func (s *transactionService) CreateTransaction(ctx context.Context, params db.CreateTransactionParams) (db.Transaction, error) {
	tx, err := s.queries.CreateTransaction(ctx, params)
	if err != nil {
		return db.Transaction{}, utils.ServiceError{
			Message:     "Failed to create transaction",
			ServiceName: "Transaction",
			ErrorMsg:    err,
		}
	}
	return tx, nil
}

func (s *transactionService) ListTransactionsByUserID(ctx context.Context, userID int32) ([]db.Transaction, error) {
	txs, err := s.queries.ListTransactionsByUserID(ctx, userID)
	if err != nil {
		return nil, utils.ServiceError{
			Message:     "Failed to list transactions",
			ServiceName: "Transaction",
			ErrorMsg:    err,
		}
	}
	return txs, nil
}
