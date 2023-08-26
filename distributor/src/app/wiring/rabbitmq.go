package wiring

import "github.com/mjedari/sternx-project/distributor/domain/contracts"

func (w *Wire) GetBroker() contracts.IBroker {
	return w.broker
}

func (w *Wire) SetNewRabbitMQInstance() error {
	return nil
}
