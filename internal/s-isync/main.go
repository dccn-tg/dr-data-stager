package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose *bool
	nworkers    *int
	configFile  *string
	drUser      *string
	drPass      *string
	stagerUser  *string
	srcPath     string
	dstPath     string
)

func init() {
	//optsConfig = flag.String("c", "config.yml", "set the `path` of the configuration file")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	nworkers = flag.Int("p", 4, "`number` of global concurrent workers")
	configFile = flag.String("c", os.Getenv("STAGER_WORKER_CONFIG"), "configurateion file `path`")
	drUser = flag.String("druser", "", "(R)DR data-access `username`")
	drPass = flag.String("drpass", "", "(R)DR data-access `password`")
	stagerUser = flag.String("user", "", "stager service local `username`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/s-isync.log",
		FileLevel:         log.Info,
	}

	if *optsVerbose {
		cfg.ConsoleLevel = log.Debug
		cfg.FileLevel = log.Debug
	}

	// initialize logger
	log.NewLogger(cfg, log.InstanceZapLogger)

	// check if both source and destination paths are provided
	args := flag.Args()
	if len(args) != 2 {
		usage()
		log.Fatalf("invalid source and desitnation paths.")
	}

	srcPath = args[0]
	dstPath = args[1]
}

func usage() {
	fmt.Printf("\ns-isync of data stager\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS] <source path> <destination path>\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancelFunc()
	}()

	// load global configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("fail to load configuration: %s", *configFile)
	}

	if err := doSomthing(ctx, cfg); err != nil {
		log.Errorf("%s\n", err)
	}

}

func doSomthing(ctx context.Context, cfg config.Configuration) error {

	log.Infof("%+v\n", cfg)

	log.Infof("%s --> %s\n", srcPath, dstPath)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
