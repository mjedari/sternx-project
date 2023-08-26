package contracts

import "context"

type IBroker interface {
	Close() error
	ProduceOnQueue(ctx context.Context, queue string, message []byte) error
	Consume(ctx context.Context, queueName, key string, apply func(ctx context.Context, message []byte) error) error
}
