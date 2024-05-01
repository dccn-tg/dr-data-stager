package main

// a simple and naive server directly using Asynqmon

import (
	"fmt"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"

	"flag"

	"os"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	optsVerbose *bool
	optsPort    *int
	redisURL    *string
	//configFile  *string
)

func usage() {
	fmt.Printf("\nStager admin interface based on Asynqmon: \n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func init() {
	optsVerbose = flag.Bool("v", false, "print debug messages")
	optsPort = flag.Int("p", 8080, "specify the service `port` number")
	redisURL = flag.String("r", "redis://redis:6379", "redis service `address`")
	//configFile = flag.String("c", os.Getenv("STAGER_ADMIN_CONFIG"), "configurateion file `path`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/admin.log",
		FileLevel:         log.Info,
	}

	if *optsVerbose {
		cfg.ConsoleLevel = log.Debug
		cfg.FileLevel = log.Debug
	}

	// initialize logger
	log.NewLogger(cfg, log.InstanceZapLogger)
}

func main() {

	// parse the redis URL to redis connection options
	redisOpts, err := asynq.ParseRedisURI(*redisURL)
	if err != nil {
		log.Fatalf("cannot parse redis URL: %s", err)
	}

	h := asynqmon.New(asynqmon.Options{
		RedisConnOpt: redisOpts,
	})

	// Note: We need the tailing slash when using net/http.ServeMux.
	http.Handle("/", h)

	log.Fatalf("%s\n", http.ListenAndServe(fmt.Sprintf(":%d", *optsPort), nil))
}
