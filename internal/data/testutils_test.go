package data

import (
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/callsamu/expenses-api/internal/testdb"
)

func SeedUsers(t *testing.T, tdb *testdb.TestDB) []*User {
	time := time.Now()

	users := []*User{
		{
			ID:        1,
			Name:      "foo",
			Email:     "foo@example.com",
			Activated: false,
			Version:   1,
			CreatedAt: time,
		},
		{
			ID:        2,
			Name:      "bar",
			Email:     "bar@example.com",
			Activated: true,
			Version:   1,
			CreatedAt: time,
		},
	}

	for _, user := range users {
		err := user.Password.Set("password")
		if err != nil {
			t.Fatal(err)
		}

		stmt := `
			INSERT INTO users (name, email, password_hash, activated, version, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		args := []interface{}{
			user.Name,
			user.Email,
			user.Password.Hash,
			user.Activated,
			user.Version,
			user.CreatedAt,
		}

		_, err = tdb.DB.Exec(stmt, args...)
		if err != nil {
			t.Fatal(err)
		}

	}

	return users
}

func SeedTokens(t *testing.T, tdb *testdb.TestDB) []*Token {
	var tokens []*Token

	targs := []struct {
		userID int64
		scope  string
		ttl    time.Duration
	}{
		{
			userID: 1,
			scope:  ScopeActivation,
			ttl:    time.Hour,
		},
		{ // Will automatically expire
			userID: 1,
			scope:  ScopeActivation,
			ttl:    0,
		},
	}

	query := `
		INSERT INTO tokens (user_id, hash, scope, expiry)
		VALUES ($1, $2, $3, $4)
	`

	for _, args := range targs {
		token, err := generateToken(args.userID, args.ttl, args.scope)
		if err != nil {
			t.Fatal(err)
		}
		tokens = append(tokens, token)
		_, err = tdb.DB.Exec(query, token.UserID, token.Hash, token.Scope, token.Expiry)
		if err != nil {
			t.Fatal(err)
		}
	}

	return tokens
}

func SeedExpenses(t *testing.T, tdb *testdb.TestDB) []*Expense {
	expenses := []*Expense{
		{
			ID:        1,
			UserID:    1,
			Recipient: "Foo Store",
			Category:  "Foos",
			Date:      time.Now(),
			Value:     money.NewFromFloat(25, money.USD),
		},
		{
			ID:        2,
			UserID:    1,
			Recipient: "Bar Store",
			Category:  "Foos",
			Date:      time.Now().AddDate(0, 1, 5),
			Value:     money.NewFromFloat(24, money.USD),
		},
		{
			ID:        3,
			UserID:    1,
			Recipient: "FooBar Store",
			Category:  "Bars",
			Date:      time.Now().AddDate(0, 1, 1),
			Value:     money.NewFromFloat(1, money.EUR),
		},
		{
			ID:        4,
			UserID:    1,
			Recipient: "FooBarXYZ Store",
			Category:  "Bars",
			Date:      time.Now(),
			Value:     money.NewFromFloat(50, money.USD),
		},
	}

	for _, expense := range expenses {
		query := `
			INSERT INTO expenses (user_id, date, recipient, description, category, amount, currency)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		args := []any{
			expense.UserID,
			expense.Date,
			expense.Recipient,
			expense.Description,
			expense.Category,
			expense.Value.Amount(),
			expense.Value.Currency().Code,
		}

		_, err := tdb.DB.Exec(query, args...)
		if err != nil {
			t.Fatal(err)
		}
	}

	return expenses
}
