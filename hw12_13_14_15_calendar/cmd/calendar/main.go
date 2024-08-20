package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	hd "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/handlers"
	rt "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/api/routing"
	cf "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
	gr "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/server/grpc"
	hp "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/server/http"
	mm "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sq "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../", "Path to calendar configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := cf.NewCalendarConfig(configFile)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}

	logg := lg.NewSLogger(config.Logger.Level)

	ctx, cancel := context.WithCancel(context.Background())

	var storage hd.AbstractStorage
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

	routes := rt.NewRouter(logg)
	routes.AddRoutes(storage)

	httpServer := hp.NewHTTPServer(logg)
	grpcServer := gr.NewGrpcServer(storage, logg)

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
		errGr := grpcServer.Stop()
		if errGr != nil {
			logg.Error("failed to stop grpc-server", "err", err.Error())
		}
		errHTTP := httpServer.Stop(ctxtime)
		if errHTTP != nil {
			logg.Error("failed to stop http-server", "err", err.Error())
		}
		if errGr != nil || errHTTP != nil {
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
		errCh <- httpServer.Start(ctx, &config, *routes)
	}()
	go func() {
		errCh <- grpcServer.Start(&config)
	}()

	select {
	case err := <-errCh:
		panic(err)
	case <-sigs:
	}
}

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
