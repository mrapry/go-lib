package dil_service

import (
	"context"

	"github.com/mrapry/go-lib/golibshared"
)

type ServiceDil interface {
	GetDataUnit(ctx context.Context, code string) <-chan golibshared.Result
}
