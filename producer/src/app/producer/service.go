package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mjedari/sternx-project/producer/app/configs"
	"github.com/mjedari/sternx-project/producer/domain/contracts"
	"github.com/mjedari/sternx-project/producer/domain/tasks"
	"github.com/mjedari/sternx-project/producer/infra/utils"
	"github.com/sirupsen/logrus"
	"time"
)

type Service struct {
	broker contracts.IBroker
	config configs.Producer
}

func NewService(broker contracts.IBroker, config configs.Producer) *Service {
	return &Service{broker: broker, config: config}
}

func (s *Service) Run(ctx context.Context) {
	// start generating and producing
	fmt.Println("Producer started...")
	ticker := time.NewTicker(s.config.GetInterval())
	totalNumberOfTasks := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 1. generate task
			generatedTasks := tasks.GenerateTasks(s.config.NumberOfTasks)

			for _, task := range generatedTasks {

				byteTask, err := json.Marshal(task)
				if err != nil {
					logrus.Errorf("error on marshalig task: %v", err)
					continue
				}

				// todo: change the retry system
				_, err = utils.Retry(func(ctx context.Context) (any, error) {
					// produce messages on tasks-exchange with tasks-routing-key in direct type
					err = s.broker.Produce(ctx, configs.Config.Rabbit.RoutingKey, byteTask)
					if err != nil {
						return nil, err
					}

					return nil, nil
				}, 10, time.Second*2)(ctx)

				if err != nil {
					logrus.Errorf("error on producing task #%v", len(task.ID))
				}

			}

			totalNumberOfTasks += len(generatedTasks)
			fmt.Printf("published #%v tasks succssfully. \n", totalNumberOfTasks)
		}
	}
}
