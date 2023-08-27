package contracts

type IMonitoring interface {
	ReceivedTasks(queue string)
	DoneTasks(queue string)
}
