package contracts

import "context"

type IBroker interface {
	Close() error
	Consume(ctx context.Context, queueName, key string, apply func(ctx context.Context, message []byte) error) error
}
