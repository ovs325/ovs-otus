package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg/rabbitmq"
	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg/sender"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../", "Path to sender configuration file")
}

func main() {
	flag.Parse()

	cfgSn, err := sender.NewSenderConfig(configFile)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}

	cfgCl, err := config.NewCalendarConfig(cfgSn.RabbitConfigPath)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}

	logger := lg.NewSLogger(cfgCl.Logger.Level)

	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfgCl.Rabbit.GetRabbitDSN(), cfgCl.Rabbit.Queue)
	if err != nil {
		logger.Error("Error initializing RabbitMQ", "error", err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		logger.Info("Closing Sender gracefully...")
		cancel()
		rabbitMQ.Close()
		if err := recover(); err != nil {
			log.Println("Panic:", err)
		}
		fmt.Println("Sender has closed!!")
	}()

	snd, err := sender.NewSender(rabbitMQ, logger)
	if err != nil {
		panic(err)
	}

	errCh := make(chan error)

	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	go func() {
		errCh <- snd.Start(ctx)
	}()
	select {
	case err := <-errCh:
		logger.Error("Sender sent an error", "error", err.Error())
		panic(err)
	case <-sigs:
		logger.Info("Sender has received a signal")
	}
}
