package tasks

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Task struct {
	ID      uuid.UUID
	Name    string
	Payload string
}

func NewTask(name string, payload string) Task {
	id, err := uuid.NewUUID()
	if err != nil {
		logrus.Errorf("can not create task: %v", err)
		return Task{}
	}
	return Task{ID: id, Name: name, Payload: payload}
}

func GenerateTasks(number int) []Task {

	var tasks []Task
	for i := 1; i <= number; i++ {
		newTask := NewTask(fmt.Sprintf("task-%v", i), "This is payload of task")
		tasks = append(tasks, newTask)
	}

	return tasks
}
