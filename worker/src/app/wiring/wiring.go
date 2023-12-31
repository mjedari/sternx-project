package wiring

import (
	"github.com/mjedari/sternx-project/worker/app/configs"
	"github.com/mjedari/sternx-project/worker/domain/contracts"
	"github.com/sirupsen/logrus"
)

var Wiring *Wire

type Wire struct {
	broker     contracts.IBroker
	monitoring contracts.IMonitoring
	Configs    configs.Configuration
}

func NewWire(broker contracts.IBroker, monitoring contracts.IMonitoring, configs configs.Configuration) *Wire {
	return &Wire{broker: broker, monitoring: monitoring, Configs: configs}
}

func (w *Wire) ShutdownServices() {
	err := w.broker.Close()
	if err != nil {
		logrus.Errorf("error on closing broker: %v", err)
	}

	// todo: shutdown healer
}
