package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	"github.com/hibiken/asynq"

	"github.com/Donders-Institute/dr-data-stager/pkg/tasks"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	redisURL    *string
	nworkers    *int
	configFile  *string
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	nworkers = flag.Int("p", 4, "`number` of global concurrent workers")
	redisURL = flag.String("r", "redis://redis:6379", "redis service `address`")
	configFile = flag.String("c", os.Getenv("STAGER_WORKER_CONFIG"), "configurateion file `path`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/worker.log",
		FileLevel:         log.Info,
	}

	if *optsVerbose {
		cfg.ConsoleLevel = log.Debug
		cfg.FileLevel = log.Debug
	}

	// initialize logger
	log.NewLogger(cfg, log.InstanceZapLogger)
}

func usage() {
	fmt.Printf("\nworker of data stager\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	// load global configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("fail to load configuration: %s", *configFile)
	}

	// redis client instance for notifying cache update
	redisOpts, err := asynq.ParseRedisURI(*redisURL)
	if err != nil {
		log.Fatalf("%s", err)
	}

	srv := asynq.NewServer(
		redisOpts,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: *nworkers,
			// Optionally specify multiple queues with different priority.
			Queues: tasks.StagerQueues,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(n*30) * time.Second
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.Handle(tasks.TypeStager, tasks.NewStager(cfg))
	mux.Handle(tasks.TypeEmailer, tasks.NewEmailer(cfg))
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
