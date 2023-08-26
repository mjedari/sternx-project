package healer

import (
	"context"
	"fmt"
	"github.com/mjedari/sternx-project/distributor/app/configs"
	"github.com/mjedari/sternx-project/distributor/domain/contracts"
	"time"
)

type Healer struct {
	Providers []contracts.IProvider
	config    configs.Healer
}

func NewHealerService(providers []contracts.IProvider, config configs.Healer) *Healer {
	return &Healer{
		Providers: providers,
		config:    config,
	}
}

func (h *Healer) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(h.config.PingInterval))
	go func() {
		defer func() {
			ticker.Stop()
			fmt.Println("\nClosing healer...")
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, provider := range h.Providers {
					if err := provider.CheckHealth(ctx); err != nil {
						if err := provider.ResetConnection(ctx); err != nil {
							// todo: log the error
						}
					}
				}
			}
		}
	}()
}
