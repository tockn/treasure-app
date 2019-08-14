package repository

import "context"

type Tx interface {
	TxHandler(ctx context.Context, fn func(ctx context.Context) error) error
}
