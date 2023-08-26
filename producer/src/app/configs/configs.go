package configs

import (
	"fmt"
	"time"
)

var Config Configuration

type Server struct {
	Host string
	Port string
}

type Rabbit struct {
	Host         string
	Port         string
	User         string
	Pass         string
	ExchangeName string
	RoutingKey   string
}

type Healer struct {
	RetryTimes   int
	RetryDelay   int
	PingInterval int
}

type Producer struct {
	Interval      int
	NumberOfTasks int
}

type Configuration struct {
	Server   Server
	Rabbit   Rabbit
	Healer   Healer
	Producer Producer
}

func (c Configuration) GetHealerConfig() Healer {
	return c.Healer
}

func (c Rabbit) GetURL() string {
	return fmt.Sprintf("amqp://%v:%v@%v:%v/", c.User, c.Pass, c.Host, c.Port)
}

func (c Healer) GetRetryDelay() time.Duration {
	return time.Second * time.Duration(c.RetryDelay)
}

func (c Healer) GetRetryTimes() int {
	return c.RetryTimes
}

func (c Producer) GetInterval() time.Duration {
	return time.Second * time.Duration(c.Interval)
}
