package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	ppath "github.com/Donders-Institute/dr-data-stager/pkg/path"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose bool   = false
	nworkers    int    = 4
	taskID      string = "0000-0000-0000-0000"
	configFile  string = os.Getenv("STAGER_WORKER_CONFIG")
	drUser      string = "stager@ru.nl"
	drPass      string
	stagerUser  string = "root"
	srcPath     string
	dstPath     string
)

func init() {
	flag.BoolVar(&optsVerbose, "v", optsVerbose, "print debug messages")
	flag.IntVar(&nworkers, "p", nworkers, "`number` of global concurrent workers")
	flag.StringVar(&configFile, "c", configFile, "configurateion file `path`")
	flag.StringVar(&taskID, "task", taskID, "stager task `id`")
	flag.StringVar(&drUser, "druser", drUser, "(R)DR data-access `username`")
	flag.StringVar(&drPass, "drpass", drPass, "(R)DR data-access `password`")
	flag.StringVar(&stagerUser, "user", stagerUser, "stager service local `username`")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     false,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    false,
		FileLocation:      "log/s-isync.log",
		FileLevel:         log.Info,
	}

	if optsVerbose {
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
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("fail to load configuration: %s", configFile)
	}

	if err := doSomthing(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

}

func doSomthing(ctx context.Context, cfg config.Configuration) error {

	log.Infof("[%s] [%s,%s] %s --> %s\n", taskID, stagerUser, drUser, srcPath, dstPath)

	// logic of performing data transfer.
	srcPathInfo, err := ppath.GetPathInfo(ctx, srcPath)
	if err != nil {
		return fmt.Errorf("invalid source url: %s", err)
	}

	total := srcPathInfo.CountFiles(ctx)
	nsuccess := 0
	nfailure := 0

	fmt.Printf("%d,%d,%d\n", total, nsuccess, nfailure)

	dstPathInfo, _ := ppath.GetPathInfo(ctx, dstPath)
	log.Debugf("[%s] srcPathInfo: %+v", taskID, srcPathInfo)
	log.Debugf("[%s] dstPathInfo: %+v", taskID, dstPathInfo)

	processed := ppath.ScanAndSync(ctx, cfg, srcPathInfo, dstPathInfo, nworkers)

	for {
		select {
		case e, more := <-processed: // handle the output of a processed file

			if !more {
				log.Debugf("[%s] finished", taskID)
				return nil
			}

			// handle the error
			if e.Error != nil { // something went wrong
				nfailure++
				fmt.Printf("%d,%d,%d\n", total, nsuccess, nfailure)
				return fmt.Errorf("%s: %s", e.File, e.Error.Error())
			}

			// increase the counter by 1, and update the queue data
			nsuccess++
			fmt.Printf("%d,%d,%d\n", total, nsuccess, nfailure)

		case <-ctx.Done():
			// receive abort signal from parent context
			return fmt.Errorf("aborting task due to context signal")
		}
	}
}
