package data

import (
	"testing"
	"time"

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

func TestExpensesModelGetsExpenses(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()

	SeedUsers(t, tdb)
	SeedExpenses(t, tdb)

	model := ExpenseModel{DB: tdb.DB}

	cases := []struct {
		name      string
		IDs       []int64
		recipient string
		month     time.Time
		category  string
		sort      string
		page      int
		pagesize  int
	}{
		{
			name:     "gets all expenses",
			IDs:      []int64{1, 2, 3, 4},
			month:    time.Unix(0, 0),
			sort:     "id",
			page:     1,
			pagesize: 10,
		},
		{
			name:      "searches expenses by recipient",
			IDs:       []int64{3, 4},
			recipient: "%FooBar%",
			month:     time.Unix(0, 0),
			sort:      "id",
			page:      1,
			pagesize:  10,
		},
		{
			name:     "searches expenses by month",
			IDs:      []int64{2, 3},
			month:    time.Now().AddDate(0, 1, 0),
			sort:     "id",
			page:     1,
			pagesize: 10,
		},
		{
			name:     "sorts expenses by value",
			IDs:      []int64{4, 1, 2, 3},
			month:    time.Unix(0, 0),
			sort:     "-amount",
			page:     1,
			pagesize: 10,
		},
		{
			name:     "paginates expenses",
			IDs:      []int64{3, 4},
			month:    time.Unix(0, 0),
			sort:     "id",
			page:     2,
			pagesize: 2,
		},
	}

	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			filters := Filters{
				Sort:         ts.sort,
				SortSafelist: []string{ts.sort},
				Page:         ts.page,
				PageSize:     ts.pagesize,
			}

			expenses, err := model.GetAll(ts.recipient, ts.month, filters)
			if err != nil {
				t.Fatal(err)
			}

			require.Equal(t, len(ts.IDs), len(expenses), "should have same number of results")

			for i := range ts.IDs {
				assert.Equal(t, ts.IDs[i], expenses[i].ID)
			}

		})
	}
}