package workers

import (
	"fmt"
	"github.com/mjedari/sternx-project/distributor/app/configs"
	"github.com/mjedari/sternx-project/distributor/infra/utils"
)

type Worker struct {
	ID    int
	Queue string
}

func NewWorker(ID int, queue string) *Worker {
	return &Worker{ID: ID, Queue: queue}
}

type WorkerPool struct {
	pool     []*Worker
	strategy *utils.RoundRobin
}

func NewWorkerPool(config configs.Distributor) *WorkerPool {
	var workerPool WorkerPool

	var input []int
	for i := 0; i < config.Workers; i++ {
		input = append(input, i)
	}

	strategy := utils.NewRoundRobin(input)
	workerPool.strategy = strategy

	for i := 0; i < configs.Config.Distributor.Workers; i++ {
		newWorker := NewWorker(i, fmt.Sprintf("worker-queue-%v", i))
		workerPool.pool = append(workerPool.pool, newWorker)
	}
	return &workerPool
}

func (p *WorkerPool) GetNextWorker() *Worker {
	nextWorkerID := p.strategy.Next()

	return p.pool[nextWorkerID]

}
