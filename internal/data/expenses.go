package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/Rhymond/go-money"
)

type Expense struct {
	ID          int64
	UserID      int64
	Date        time.Time
	Recipient   string
	Description string
	Category    string
	Value       *money.Money
	Version     int
}

type ExpenseModel struct {
	DB *sql.DB
}

func (m ExpenseModel) Insert(expense *Expense) error {
	query := `
		INSERT INTO expenses (user_id, recipient, description, category, amount, currency)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, date, version
	`

	args := []any{
		expense.UserID,
		expense.Recipient,
		expense.Description,
		expense.Category,
		expense.Value.Amount(),
		expense.Value.Currency().Code,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&expense.ID,
		&expense.Date,
		&expense.Version,
	)

	if err != nil {
		return err
	}

	return nil
}
