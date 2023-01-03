package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenIsGenerated(t *testing.T) {
	token, err := generateToken(1, time.Second, "")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int64(1), token.UserID)

	time := time.Now()
	assert.Equal(t, time.Minute(), token.Expiry.Minute())
	assert.Equal(t, token.Scope, "")

	assert.True(t, len(token.Plaintext) == 26, "token plaintext size must be 26 bytes long")
	assert.NotNil(t, token.Hash)
}
