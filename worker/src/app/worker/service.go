package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mjedari/sternx-project/worker/app/configs"
	"github.com/mjedari/sternx-project/worker/domain/contracts"
	"github.com/mjedari/sternx-project/worker/domain/tasks"
	"github.com/sirupsen/logrus"
)

type Service struct {
	broker contracts.IBroker
	config configs.Worker
}

func NewService(broker contracts.IBroker, config configs.Worker) *Service {
	return &Service{broker: broker, config: config}
}

func (s *Service) Run(ctx context.Context) {
	// start generating and producing
	fmt.Printf("Worker started on queue: %v \n", s.config.QueueName)

	err := s.broker.Consume(ctx, s.config.QueueName, "", s.handleMessage)
	if err != nil {
		logrus.Fatalf("error on consuming")
	}

}

func (s *Service) handleMessage(ctx context.Context, message []byte) error {
	fmt.Printf("message received on queue: %v \n", s.config.QueueName)

	var task tasks.Task
	err := json.Unmarshal(message, &task)
	if err != nil {
		return err
	}

	fmt.Println(task)

	return nil
}
