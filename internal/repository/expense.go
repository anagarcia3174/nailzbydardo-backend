package repository

import (
	"context"

	"fmt"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/model"
	"time"
)

type ExpenseRepository struct {
	q *sqlc.Queries
}

func NewExpenseRepository(q *sqlc.Queries) *ExpenseRepository {
	return &ExpenseRepository{q: q}
}

func (r *ExpenseRepository) CreateExpense(ctx context.Context, name string, price int64, datePurchased time.Time, receiptURL *string) (model.Expense, error) {
	createExpenseParams := sqlc.CreateExpenseParams{
		ExpenseName:   name,
		Price:         centsToPgNumeric(price),
		DatePurchased: timeToPgDate(datePurchased),
		ReceiptUrl:    stringPtrToPgText(receiptURL),
	}

	e, err := r.q.CreateExpense(ctx, createExpenseParams)

	if err != nil {
		return model.Expense{}, fmt.Errorf("error creating expense: %w", err)
	}
	id, err := pgUUIDToString(e.ID)
	if err != nil {
		return model.Expense{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	expense := model.Expense{
		ID:            id,
		ExpenseName:   e.ExpenseName,
		Price:         pgNumericToCents(e.Price),
		DatePurchased: pgDateToTime(e.DatePurchased),
		ReceiptURL:    pgTextToStringPtr(e.ReceiptUrl),
		CreatedAt:     pgTimestamptzToTime(e.CreatedAt),
	}
	return expense, nil
}

func (r *ExpenseRepository) ListExpenses(ctx context.Context) ([]model.Expense, error) {
	expenseRows, err := r.q.ListExpenses(ctx)
	if err != nil {
		return []model.Expense{}, fmt.Errorf("error getting expenses: %w", err)
	}
	expenses := []model.Expense{}
	for _, e := range expenseRows {
		id, err := pgUUIDToString(e.ID)
		if err != nil {
			return []model.Expense{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		expense := model.Expense{
			ID:            id,
			ExpenseName:   e.ExpenseName,
			Price:         pgNumericToCents(e.Price),
			DatePurchased: pgDateToTime(e.DatePurchased),
			ReceiptURL:    pgTextToStringPtr(e.ReceiptUrl),
			CreatedAt:     pgTimestamptzToTime(e.CreatedAt),
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (r *ExpenseRepository) DeleteExpense(ctx context.Context, id string) error {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return ErrInvalidID
	}
	err = r.q.DeleteExpense(ctx, pgID)
	if err != nil {
		return fmt.Errorf("error deleting expense: %w",  err)
	}
	return nil
}

func (r *ExpenseRepository) GetExpensesForPeriod(ctx context.Context, dateOne time.Time, dateTwo time.Time) (int64, error) {
	getExpensesForPeriodParams := sqlc.GetExpensesForPeriodParams{
		DatePurchased: timeToPgDate(dateOne),
		DatePurchased_2: timeToPgDate(dateTwo),
	}
	t, err := r.q.GetExpensesForPeriod(ctx, getExpensesForPeriodParams)

	if err != nil {
		return 0,  fmt.Errorf("error getting expenses for period: %w", err)
	}
	
	return pgNumericToCents(t), nil
}