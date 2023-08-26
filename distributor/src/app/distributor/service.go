package distributor

import (
	"context"
	"fmt"
	"github.com/mjedari/sternx-project/distributor/app/configs"
	"github.com/mjedari/sternx-project/distributor/domain/contracts"
	"github.com/mjedari/sternx-project/distributor/domain/workers"
	"github.com/sirupsen/logrus"
)

type Service struct {
	broker  contracts.IBroker
	workers *workers.WorkerPool
	config  configs.Distributor
}

func NewService(broker contracts.IBroker, workers *workers.WorkerPool, config configs.Distributor) *Service {
	return &Service{broker: broker, workers: workers, config: config}
}

func (s *Service) Run(ctx context.Context) {
	// start generating and producing
	fmt.Println("Distributor started...")
	fmt.Println("Workers:", s.workers.GetNextWorker())

	err := s.broker.Consume(ctx, "", "", s.distributeMessages)
	if err != nil {
		logrus.Errorf("error on consuming: %v", err)
	}

}

func (s *Service) distributeMessages(ctx context.Context, message []byte) error {

	// distribute messages
	worker := s.workers.GetNextWorker()
	//worker.Queue
	queueName := worker.Queue

	fmt.Println("produced on message on ", queueName)

	// publish the message on queue
	err := s.broker.ProduceOnQueue(ctx, queueName, message)
	if err != nil {
		// todo: handel this
		return err
	}

	return nil
}
