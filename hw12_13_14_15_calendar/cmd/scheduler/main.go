package scheduler

// import (
//     "context"
//     "log"
//     "time"

//     "gopkg.in/yaml.v2"
//     "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/rabbitmq"
//     "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/scheduler"
//     "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/sender"
// 	cf "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/config"
// )

// func main() {
// 	...
// 	var cfg cf.Config
// 	rabbitDSN := fmt.Spintf(
// 		"amqp://%s:%s@%s:%v/",
// 		cfg.Rabbit.User,
// 		cfg.Rabbit.Password,
// 		cfg.Rabbit.Host,
// 		cfg.Rabbir.Port,
// 	)
//     // Подключение к RabbitMQ
//     rabbitMQ, err := rabbitmq.NewRabbitMQ(rabbitDSN, cfg.RabbitMQ.Queue)
//     if err != nil {
//         log.Fatalf("Error initializing RabbitMQ: %v", err)
//     }
//     defer rabbitMQ.Close()

//     // Создание контекста
//     ctx, cancel := context.WithCancel(context.Background())
//     defer cancel()

//     // Запуск планировщика
//     sched := scheduler.NewScheduler(db, rabbitMQ, parseInterval(cfg.Scheduler.Interval))
//     go sched.Start(ctx)

//     // Запуск рассыльщика
//     snd := sender.NewSender(rabbitMQ)
//     go snd.Start(ctx)

//     // Ожидание завершения процессов
//     select {}
// }

// func parseInterval(interval string) time.Duration {
//     d, err := time.ParseDuration(interval)
//     if err != nil {
//         log.Fatalf("Error parsing interval: %v", err)
//     }
//     return d
// }

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

	select {
	case err := <-errCh:
		panic(err)
	case <-sigs:
	}
}
