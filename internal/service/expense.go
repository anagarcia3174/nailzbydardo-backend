package service

import (
	"context"
	"nailzbydardo/internal/model"
	"nailzbydardo/internal/repository"
	"strings"
	"time"
)

// ExpenseService struct, holding *repository.ExpenseRepository
type ExpenseService struct {
	expenseRepo *repository.ExpenseRepository
}

func NewExpenseService(expenseRepo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, name string, price int64, datePurchased time.Time, receiptURL *string) (model.Expense, error) {
	if strings.TrimSpace(name) == "" {
		return model.Expense{}, ErrInvalidInput
	}
	return s.expenseRepo.CreateExpense(ctx, name, price, datePurchased, receiptURL)
}

func (s *ExpenseService) ListExpenses(ctx context.Context) ([]model.Expense, error) {
	return s.expenseRepo.ListExpenses(ctx)
}

func (s *ExpenseService)  DeleteExpense(ctx context.Context, id string) error {
	return s.expenseRepo.DeleteExpense(ctx, id)
}

func (s *ExpenseService) GetExpensesForPeriod(ctx context.Context, dateOne time.Time, dateTwo time.Time) (int64, error) {
	return s.expenseRepo.GetExpensesForPeriod(ctx, dateOne, dateTwo)
}