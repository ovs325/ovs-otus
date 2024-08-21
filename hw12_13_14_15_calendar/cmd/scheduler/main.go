package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
	mm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sq "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg/rabbitmq"
	sh "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg/scheduler"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../", "Path to scheduller configuration file")
}

func main() {
	flag.Parse()

	cfgSh, err := sh.NewShedullerConfig(configFile)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}

	cfgCl, err := config.NewCalendarConfig(cfgSh.RabbitConfigPath)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}

	logger := lg.NewSLogger(cfgCl.Logger.Level)

	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfgCl.Rabbit.GetRabbitDSN(), cfgCl.Rabbit.Queue)
	if err != nil {
		log.Fatalf("Error initializing RabbitMQ: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		logger.Info("Closing Scheduller gracefully...")
		cancel()
		rabbitMQ.Close()
		if err := recover(); err != nil {
			log.Println("Panic:", err)
		}
		fmt.Println("Scheduller has closed!!")
	}()

	var storage hd.AbstractStorage
	var errStorage error
	if cfgCl.DB.IsPostgres {
		storage, errStorage = sq.NewPgRepo(ctx, &cfgCl, logger)
	} else {
		storage, errStorage = mm.NewMemRepo()
	}
	if errStorage != nil {
		logger.Error("data storage initialization error", "err", errStorage.Error())
		return
	}

	scheduler, err := sh.NewScheduler(storage.(sh.Repo), rabbitMQ, cfgSh.Interval)
	if err != nil {
		panic(err)
	}
	logger.Info("Scheduller:", "interval", cfgSh.Interval)

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
		errCh <- scheduler.Start(ctx)
	}()

	logger.Info("Scheduller is start")
	select {
	case err := <-errCh:
		logger.Error("Scheduller sent an error")
		panic(err)
	case <-sigs:
		logger.Info("Scheduller has received a signal")
	}
}
