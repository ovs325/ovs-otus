package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rt "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/routing"
	bl "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/business_logic"
	cf "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
	hp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/server/http"
	mm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sq "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := cf.NewConfig(configFile)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}

	logg := lg.NewSLogger(config.Logger.Level)

	ctx, cancel := context.WithCancel(context.Background())

	var storage bl.AbstractStorage
	var errStorage error
	if config.DB.IsPostgres {
		storage, errStorage = sq.NewPgRepo(ctx, &config, logg)
	} else {
		storage, errStorage = mm.NewMemRepo()
	}
	if errStorage != nil {
		logg.Error("data storage initialization error", "err", errStorage.Error())
		return
	}

	logic := bl.NewBusinessLogic(storage)

	routes := rt.NewRouter(logg)
	routes.AddRoutes(logic)

	server := hp.NewServer(logg)

	// init graceful shutdown.
	defer func() {
		logg.Info("Closing microservice gracefully...")
		storage.Close()
		cancel()
		if err := recover(); err != nil {
			log.Println("Panic:", err)
		}
		ctxtime, canceltime := context.WithTimeout(context.Background(), time.Second*3)
		defer canceltime()
		if err := server.Stop(ctxtime); err != nil {
			logg.Error("failed to stop http server", "err", err.Error())
			os.Exit(1)
		}
		fmt.Println("Microservice has closed!!")
	}()

	// start server.
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
		errCh <- server.Start(ctx, &config, *routes)
	}()

	select {
	case err := <-errCh:
		panic(err)
	case <-sigs:
	}
}
