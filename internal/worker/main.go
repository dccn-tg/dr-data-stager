package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/logging"
	"github.com/thoas/bokchoy/middleware"

	hapi "github.com/Donders-Institute/dr-data-stager/internal/api-server/handler"
	hworker "github.com/Donders-Institute/dr-data-stager/internal/worker/handler"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
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
	// initiate blochy queue for setting project roles
	var logger logging.Logger

	// redis client instance for notifying cache update
	redisOpts, err := redis.ParseURL(*redisURL)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ctx := context.Background()

	// connect to bokchoy task scheduler.
	// see https://github.com/thoas/bokchoy
	bok, err := bokchoy.New(ctx,
		bokchoy.Config{
			Broker: bokchoy.BrokerConfig{
				Type: "redis",
				Redis: bokchoy.RedisConfig{
					Type: "client",
					Client: bokchoy.RedisClientConfig{
						Addr: redisOpts.Addr,
					},
				},
			},
		},
		bokchoy.WithLogger(logger),
		bokchoy.WithTTL(3*24*time.Hour),
		bokchoy.WithConcurrency(*nworkers),
	)

	if err != nil {
		log.Errorf("cannot connect to db: %s", err)
		os.Exit(1)
	}

	//NOTE: we don't use `middleware.Timout` because that will
	//      immediately cancel the job.  What we want is to fail
	//      the job followed by another retry.  Thus  we handle
	//      timout within the main job handler.
	bok.Use(middleware.Recoverer)
	bok.Use(middleware.DefaultLogger)

	// add handler to handle tasks in the queue
	queue := bok.Queue(hapi.StagerJobQueueName)

	// cache event when job is started
	queue.OnStartFunc(func(r *bokchoy.Request) error {

		// we need to registry `Task.MaxRetries` as context value
		// so that it can be used in the event handler when task
		// is failed.
		*r = *r.WithContext(context.WithValue(r.Context(), "maxRetries", r.Task.MaxRetries))

		// log start of the job
		log.Infof("[%s] started, max retries: %d", r.Task.ID, r.Task.MaxRetries)
		return nil
	})

	// cache event when job is succeeded
	queue.OnSuccessFunc(func(r *bokchoy.Request) error {
		// log success of the job
		log.Infof("[%s] succeeded", r.Task.ID)
		// TODO: notify job owner
		return nil
	})

	// cache event when job is failed
	queue.OnFailureFunc(func(r *bokchoy.Request) error {

		// get remaining retries from context
		maxRetries := r.Context().Value("maxRetries")

		// log failure of the job
		log.Infof("[%s] failed, remaing attempts: %d", r.Task.ID, maxRetries)

		if maxRetries == 0 { // task failed definitely.
			// TODO: notify administrator
			// TODO: optionally notify job owner
			log.Infof("[%s] sending alert on failue", r.Task.ID)
		}

		return nil
	})

	// // cache event when job is completed
	// queue.OnCompleteFunc(func(r *bokchoy.Request) error {
	// 	switch {
	// 	case r.Task.IsStatusSucceeded():
	// 		log.Infof("[%s] completed successfully", r.Task.ID)
	// 		// TODO: send notification to owner
	// 		return nil
	// 	case r.Task.IsStatusFailed():
	// 		log.Infof("[%s] completed with failure", r.Task.ID)
	// 		// TODO: send notification to admin
	// 		return nil
	// 	case r.Task.IsStatusCanceled():
	// 		log.Infof("[%s] canceled", r.Task.ID)
	// 		return nil
	// 	default:
	// 		return nil
	// 	}
	// })

	// main job handler
	queue.Handle(
		&hworker.StagerJobRunner{
			ConfigFile: *configFile,
			Queue:      queue,
		},
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Warnf("Received signal, gracefully stopping")
			bok.Stop(ctx)
		}
	}()

	bok.Run(ctx)
}
