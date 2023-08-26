package contracts

import "context"

type IBroker interface {
	Close() error
	Produce(ctx context.Context, key string, message []byte) error
}
