package mocks

import (
	"time"

	"github.com/callsamu/expenses-api/internal/data"
)

var MockNonActivatedUser = &data.User{
	ID:        1,
	Name:      "foo",
	Email:     "foo@example.com",
	Version:   1,
	Activated: false,
	CreatedAt: time.Now(),
}

var MockActivatedUser = &data.User{
	ID:        2,
	Name:      "bar",
	Email:     "bar@example.com",
	Version:   1,
	Activated: true,
	CreatedAt: time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(user *data.User) error {

	if user.Email == MockActivatedUser.Email || user.Email == MockNonActivatedUser.Email {
		return data.ErrDuplicateEmail
	}

	user.ID = 3
	user.Version = 1
	user.CreatedAt = time.Now()

	return nil
}

func (m *UserModel) GetByEmail(email string) (*data.User, error) {
	switch email {
	case "foo@example.com":
		MockNonActivatedUser.Password.Set("mypassword")
		return MockNonActivatedUser, nil
	case "bar@example.com":
		MockActivatedUser.Password.Set("mypassword")
		return MockActivatedUser, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m *UserModel) Update(user *data.User) error {
	if user.ID == 1 && user.Email == MockActivatedUser.Email {
		return data.ErrDuplicateEmail
	} else if user.ID == 2 && user.Email == MockNonActivatedUser.Email {
		return data.ErrDuplicateEmail
	}

	if user.Version != 1 {
		return data.ErrEditConflict
	}

	user.Version += 1

	return nil
}

func (m *UserModel) GetForToken(scope string, plaintext string) (*data.User, error) {
	switch scope {
	case data.ScopeActivation:
		if plaintext == MockActivationToken.Plaintext {
			return MockNonActivatedUser, nil
		}
		return nil, data.ErrRecordNotFound
	case data.ScopeAuthentication:
		if plaintext == MockAuthenticationToken.Plaintext {
			return MockActivatedUser, nil
		}
		return nil, data.ErrRecordNotFound
	}
	return nil, data.ErrRecordNotFound
}
