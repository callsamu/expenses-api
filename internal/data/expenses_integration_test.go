package data

import (
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/callsamu/expenses-api/internal/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpenseModelInsertExpense(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()

	SeedUsers(t, tdb)

	model := ExpenseModel{DB: tdb.DB}
	expense := &Expense{
		UserID:      1,
		Recipient:   "Foo Store",
		Description: "Transaction at Foo Store",
		Category:    "Miscelaneous",
		Value:       money.NewFromFloat(1.0, money.USD),
	}

	t.Run("properly inserts expense", func(t *testing.T) {
		err := model.Insert(expense)
		if err != nil {
			t.Fatal(err)
		}

		var output struct {
			Amount   int64
			Currency string
		}

		query := "SELECT amount, currency FROM expenses WHERE id = 1"
		err = tdb.DB.QueryRow(query).Scan(&output.Amount, &output.Currency)
		require.Nil(t, err)

		assert.EqualValues(t, 1, expense.ID)
		assert.EqualValues(t, 1, expense.Version)
		assert.NotZero(t, expense.Date, "expected date not to be null")

		value := money.New(output.Amount, output.Currency)
		isEqual, err := expense.Value.Equals(value)
		if err != nil {
			t.Fatal(err)
		}
		assert.True(t, isEqual, "monetary values should be equal")

	})
}