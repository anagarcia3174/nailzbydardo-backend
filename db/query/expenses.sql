-- name: CreateExpense :one
INSERT INTO expenses (expense_name, price, date_purchased, receipt_url) VALUES ($1, $2, $3, $4) RETURNING id, expense_name, price, date_purchased, receipt_url, created_at;

-- name: ListExpenses :many
SELECT id, expense_name, price, date_purchased, receipt_url, created_at FROM expenses;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1;

-- name: GetExpensesForPeriod :one
SELECT COALESCE(SUM(price), 0)::numeric AS total FROM expenses
WHERE date_purchased >= $1 AND date_purchased < $2;