package main

import (
	"context"
	"errors"
	"strings"
)

type stringService struct{}

func (stringService) Uppercase(_ context.Context, s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(_ context.Context, s string) int {
	return len(s)
}

// ErrEmpty returned when input string is empty
var ErrEmpty = errors.New("String is empty")
