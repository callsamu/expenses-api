package data

import (
	"testing"
	"time"

	"github.com/callsamu/pfapi/internal/testdb"
)

func SeedUsers(t *testing.T, tdb *testdb.TestDB) []User {
	time := time.Now()

	users := []User{
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
		ptr := &user
		err := ptr.Password.Set("password")
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
