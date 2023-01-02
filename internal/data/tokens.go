package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	UserID    int64
	Plaintext string
	Hash      []byte
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	plaintext := base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(bytes)
	hash := sha256.Sum256([]byte(plaintext))

	token := &Token{
		UserID:    userID,
		Plaintext: plaintext,
		Hash:      hash[:],
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	return token, nil
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (user_id, hash, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	args := []interface{}{token.UserID, token.Hash, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	if err != nil {
		return err
	}

	return nil
}
