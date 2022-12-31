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
