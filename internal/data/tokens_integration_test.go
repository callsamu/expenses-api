package data

import (
	"database/sql"
	"testing"
	"time"

	"github.com/callsamu/expenses-api/internal/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenIsCreatedAndInserted(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()
	user := SeedUsers(t, tdb)[0]

	model := TokenModel{DB: tdb.DB}

	t.Run("inserts token into DB", func(t *testing.T) {
		token, err := model.New(user.ID, time.Second, ScopeActivation)

		var hash []byte
		query := "SELECT hash FROM tokens WHERE user_id = $1"
		err = tdb.DB.QueryRow(query, user.ID).Scan(&hash)
		require.Nil(t, err)

		assert.Equal(t, hash, token.Hash)
	})
}

func TestAllTokensCanBeDeletedPerUser(t *testing.T) {
	if testing.Short() {
		t.Skip("data: skipping integration test")
	}

	tdb := testdb.Open(t)
	defer tdb.Close()

	SeedUsers(t, tdb)
	tokens := SeedTokens(t, tdb)

	var tokensForUser []*Token
	for _, token := range tokens {
		if token.UserID == 1 {
			tokensForUser = append(tokensForUser, token)
		}
	}

	model := TokenModel{DB: tdb.DB}

	t.Run("deletes all tokens per user and activation scope", func(t *testing.T) {
		model.DeleteAllForUser(ScopeActivation, 1)

		for _, token := range tokensForUser {
			var result []byte
			query := `SELECT * FROM tokens WHERE hash = $1`
			err := tdb.DB.QueryRow(query, token.Hash).Scan(&result)
			require.ErrorIs(t, err, sql.ErrNoRows)
		}
	})
}
