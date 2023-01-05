package mocks

import (
	"crypto/sha256"
	"time"

	"github.com/callsamu/expenses-api/internal/data"
	"github.com/stretchr/testify/mock"
)

var MockActivationToken = &data.Token{
	UserID:    1,
	Plaintext: "abcabcabcabcabcabcabcabc01",
	Expiry:    time.Now().Add(time.Hour),
	Scope:     data.ScopeActivation,
}

//var MockAuthenticationToken = &data.Token{
//	UserID:    1,
//	Plaintext: "abcabcabcabcabcabcabcabc02",
//	Expiry:    time.Now().Add(time.Hour),
//	Scope:     data.,
//}

type TokenModel struct {
	mock.Mock
}

func (m *TokenModel) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	plaintext := "abcabcabcabcabcabcabcabc00"
	hash := sha256.Sum256([]byte(plaintext))

	token := &data.Token{
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Plaintext: plaintext,
		Hash:      hash[:],
	}

	return token, nil
}

func (m *TokenModel) Insert(token *data.Token) error {
	return nil
}

func (m *TokenModel) DeleteAllForUser(scope string, userID int64) error {
	return nil
}
