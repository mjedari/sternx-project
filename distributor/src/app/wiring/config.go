package wiring

import (
	"net"
)

func (w *Wire) GetRedisUrl() string {
	return net.JoinHostPort(w.Configs.Server.Host, w.Configs.Server.Port)
}

func (w *Wire) GetServerConfig() string {
	return net.JoinHostPort(w.Configs.Server.Host, w.Configs.Server.Port)
}
