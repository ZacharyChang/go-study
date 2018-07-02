package main

import "context"

type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}
