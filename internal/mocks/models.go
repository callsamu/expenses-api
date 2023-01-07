package mocks

import "github.com/callsamu/expenses-api/internal/data"

func NewModels() data.Models {
	return data.Models{
		Users:  &UserModel{},
		Tokens: &TokenModel{},
	}
}
