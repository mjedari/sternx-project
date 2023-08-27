package wiring

import "github.com/mjedari/sternx-project/worker/domain/contracts"

func (w *Wire) GetMonitoringService() contracts.IMonitoring {
	return w.monitoring
}
