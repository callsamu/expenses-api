package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type password struct {
	Plaintext *string
	Hash      []byte
}

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  password
	Version   int
	Activated bool
	CreatedAt time.Time
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintext
	p.Hash = hash

	return nil
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version, created_at
	`

	email := strings.ToLower(user.Email)
	args := []interface{}{user.Name, email, user.Password.Hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Version, &user.CreatedAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM users
		WHERE email = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
