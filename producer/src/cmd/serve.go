package cmd

import (
	"context"
	"fmt"
	"github.com/mjedari/sternx-project/producer/app/configs"
	"github.com/mjedari/sternx-project/producer/app/producer"
	"github.com/mjedari/sternx-project/producer/app/wiring"
	"github.com/mjedari/sternx-project/producer/domain/contracts"
	"github.com/mjedari/sternx-project/producer/infra/broker"
	"github.com/mjedari/sternx-project/producer/infra/healer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serving service.",
	Long:  `Serving service.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())

		serve(ctx)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
		// shutting down services
		wiring.Wiring.ShutdownServices()

		//fmt.Println()
		//for i := 5; i > 0; i-- {
		//	time.Sleep(time.Second)
		//	fmt.Printf("\033[2K\rShutting down ... : %d", i)
		//}

		// Perform any necessary cleanup before exiting
		fmt.Println("\nService exited successfully.")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(ctx context.Context) {
	initWiring(ctx)

	service := producer.NewService(wiring.Wiring.GetBroker(), wiring.Wiring.GetProducerConfig())
	service.Run(ctx)
	// start project from here
}

func initWiring(ctx context.Context) {
	fmt.Println("configs.Config.Rabbit", configs.Config.Rabbit)
	rabbitProvider, err := broker.NewRabbitMQ(configs.Config.Rabbit)
	if err != nil {
		logrus.Fatalf("Fatal error on create redis("+configs.Config.Rabbit.Host+":"+configs.Config.Rabbit.Port+")connection: %s \n", err)
	}

	wiring.Wiring = wiring.NewWire(rabbitProvider, configs.Config)

	// init healer for services
	infraHealer := healer.NewHealerService([]contracts.IProvider{rabbitProvider}, configs.Config.GetHealerConfig())
	infraHealer.Start(ctx)

	logrus.Info("wiring initialized")
}
