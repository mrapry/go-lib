package caterpillar_service

import (
	"context"

	"github.com/mrapry/go-lib/golibshared"
)

type ServiceCaterpillar interface {
	GetDataUnit(ctx context.Context, code string) <-chan golibshared.Result
}
