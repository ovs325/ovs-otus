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

	rt "bruteforce/api/routing"
	cf "bruteforce/config"
	sr "bruteforce/internal"
	lg "bruteforce/internal/logger"
	mn "bruteforce/internal/manager"
	hp "bruteforce/internal/server/http"
)

const configPath = "./"

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", configPath, "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := cf.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}

	if config.IsTest {
		sr.SetElapsed(0)
		sr.GetNowMu.Lock()
		sr.GetNow = func() time.Time { return sr.Start.Add(sr.GetElapsed()) }
		sr.GetNowMu.Unlock()
	}

	logg := lg.NewSLogger(config.Logger.Level)

	logg.Info("I'v got Config path: ", "path", configFile)
	logg.Info("I'm looking for a configuration here: ", "path", fmt.Sprintf("%sconfig.yaml", configFile))
	logg.Info("Params:", "Host", config.HTTPServer.Host, "Port", config.HTTPServer.Port)

	ctx, cancel := context.WithCancel(context.Background())

	routes := rt.NewRouter(logg)
	mng, err := mn.NewManager(configPath)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
	}
	routes.AddRoutes(mng, config.IsTest)

	httpServer := hp.NewHTTPServer(logg)

	// init graceful shutdown.
	defer func() {
		logg.Info("Closing microservice gracefully...")
		cancel()
		if err := recover(); err != nil {
			log.Println("Panic:", err)
		}
		ctxtime, canceltime := context.WithTimeout(context.Background(), time.Second*3)
		defer canceltime()
		errHTTP := httpServer.Stop(ctxtime)
		if errHTTP != nil {
			logg.Error("failed to stop http-server", "err", err.Error())
		}
		if errHTTP != nil {
			os.Exit(1)
		}
		fmt.Println("Microservice has closed")
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
