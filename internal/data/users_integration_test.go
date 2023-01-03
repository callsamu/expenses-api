package data

import (
	"testing"

	"github.com/callsamu/expenses-api/internal/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserModelInsertsUser(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()

	model := UserModel{DB: tdb.DB}

	user := &User{
		Name:  "foo",
		Email: "foo@example.com",
	}

	err := user.Password.Set("password")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("properly inserts user", func(t *testing.T) {
		err = model.Insert(user)
		if err != nil {
			t.Fatal(err)
		}

		var count int64
		query := "SELECT COUNT(id) FROM users"
		err = tdb.DB.QueryRow(query).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}

		assert.EqualValues(t, 1, count, "expected users count to be 1")
		assert.EqualValues(t, 1, user.ID, "expected user ID to be 1")
		assert.EqualValues(t, 1, user.Version, "expected user version to be 1")
	})

	t.Run("returns ErrDuplicateEmail if user email is duplicated", func(t *testing.T) {
		err := model.Insert(user)
		assert.ErrorIs(t, err, ErrDuplicateEmail)
	})
}

func TestUserModelFindsUserByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()
	users := SeedUsers(t, tdb)
	model := UserModel{DB: tdb.DB}

	t.Run("finds the right user by email", func(t *testing.T) {
		user, err := model.GetByEmail("foo@example.com")
		if err != nil {
			t.Fatal(err)
		}

		assert.EqualValues(t, users[0].ID, user.ID, "users ID should match")
	})

	t.Run("returns ErrRecordNotFound", func(t *testing.T) {
		_, err := model.GetByEmail("foobar@example.com")
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})
}

func TestUserModelUpdatesUser(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()
	user := SeedUsers(t, tdb)[0]
	model := UserModel{DB: tdb.DB}

	t.Run("updates user", func(t *testing.T) {
		user.Name = "foobar"
		err := model.Update(&user)
		if err != nil {
			t.Fatal(err)
		}

		var name string
		query := `SELECT name FROM users WHERE id = $1`
		err = tdb.DB.QueryRow(query, user.ID).Scan(&name)

		assert.Equal(t, user.Name, name, "expected user name to be updated")
	})

	t.Run("optimistic locks", func(t *testing.T) {
		user.Version = 2
		err := model.Update(&user)
		assert.ErrorIs(t, err, ErrEditConflict)
	})

}

func TestUserModelFindsUserByToken(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()

	user := SeedUsers(t, tdb)[0]
	tokens := SeedTokens(t, tdb)

	model := UserModel{DB: tdb.DB}

	cases := []struct {
		name      string
		plaintext string
		err       error
	}{
		{
			name:      "gets right user for token",
			plaintext: tokens[0].Plaintext,
			err:       nil,
		},
		{
			name:      "ignores expired tokens",
			plaintext: tokens[1].Plaintext,
			err:       ErrRecordNotFound,
		},
	}

	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			retrievedUser, err := model.GetForToken(ScopeActivation, ts.plaintext)
			require.ErrorIs(t, err, ts.err)
			if ts.err != nil {
				return
			}
			assert.Equal(t, user.ID, retrievedUser.ID)
			assert.Equal(t, user.Email, retrievedUser.Email)
		})
	}
}
