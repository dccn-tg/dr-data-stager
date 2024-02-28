package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	"github.com/Donders-Institute/dr-data-stager/pkg/dr"
	"github.com/Donders-Institute/dr-data-stager/pkg/errors"
	ppath "github.com/Donders-Institute/dr-data-stager/pkg/path"
	"github.com/Donders-Institute/dr-data-stager/pkg/utility"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

var (
	//optsConfig  *string
	optsVerbose       bool   = false
	nworkers          int    = 4
	taskID            string = "0000-0000-0000-0000"
	logFile           string = "/opt/stager/log/s-isync.log"
	configFile        string = os.Getenv("STAGER_WORKER_CONFIG")
	drUser            string = "stager@ru.nl"
	drPass            string
	withEncryptedPass bool   = false
	rsaKey            string = "key.pem"
	srcPath           string
	dstPath           string
)

func init() {
	flag.BoolVar(&optsVerbose, "v", optsVerbose, "print debug messages")
	flag.IntVar(&nworkers, "p", nworkers, "`number` of global concurrent workers")
	flag.StringVar(&configFile, "c", configFile, "configurateion file `path`")
	flag.StringVar(&logFile, "l", logFile, "log file `path`")
	flag.StringVar(&taskID, "task", taskID, "stager task `id`")
	flag.StringVar(&drUser, "druser", drUser, "(R)DR data-access `username`")
	flag.StringVar(&drPass, "drpass", drPass, "(R)DR data-access `password`")
	flag.BoolVar(&withEncryptedPass, "e", withEncryptedPass, "use encrypted (R)DR data-access password")
	flag.StringVar(&rsaKey, "k", rsaKey, "RSA key `path` for decrypting (R)DR data-access password")

	flag.Usage = usage

	flag.Parse()

	cfg := log.Configuration{
		EnableConsole:     false,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Info,
		EnableFile:        true,
		FileJSONFormat:    false,
		FileLocation:      logFile,
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
		fmt.Fprintf(os.Stderr, "Error: insufficient arguments for source and destination")
		os.Exit(128) // invalid argument
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

	ctx, cancel := context.WithCancel(context.Background())

	if withEncryptedPass {
		encrypted, err := utility.DecryptStringWithRsaKey(drPass, rsaKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: fail to decrypt credential: %s", err)
			os.Exit(128)
		}
		drPass = *encrypted
	}

	ctx = context.WithValue(ctx, dr.KeyCredential, dr.NewCredential(drUser, drPass))

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	// load global configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: fail to load configuration: %s", configFile)
		os.Exit(128) // invalid argument
	}

	if err := run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		// TODO: define proper exit code based on the error type.
		os.Exit(err.ExitCode())
	}

}

func run(ctx context.Context, cfg config.Configuration) *errors.IsyncError {

	user, err := user.Current()
	if err != nil {
		return errors.ToIsyncError(126, err.Error())
	}

	log.Infof("[%s] [%s,%s] %s --> %s\n", taskID, user.Username, drUser, srcPath, dstPath)

	// initialize irods filesystem
	ifs, err := dr.NewFileSystem("s-isync", cfg.Dr)
	if err != nil {
		return errors.ToIsyncError(1, err.Error())
	}
	defer ifs.Release()

	ctxfs := context.WithValue(ctx, dr.KeyFilesystem, ifs)

	// logic of performing data transfer.
	srcPathInfo, err := ppath.GetPathInfo(ctxfs, srcPath)
	if err != nil {
		return errors.ToIsyncError(128, err.Error())
	}

	total := srcPathInfo.CountFiles(ctxfs)
	nsuccess := 0
	nfailure := 0

	fmt.Printf("%d,%d,%d\n", total, nsuccess, nfailure)

	dstPathInfo, _ := ppath.GetPathInfo(ctxfs, dstPath)
	log.Debugf("[%s] srcPathInfo: %+v", taskID, srcPathInfo)
	log.Debugf("[%s] dstPathInfo: %+v", taskID, dstPathInfo)

	processed := ppath.ScanAndSync(ctxfs, cfg, srcPathInfo, dstPathInfo, nworkers)

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
				return errors.ToIsyncError(1, e.Error.Error())
			}

			// increase the counter by 1, and update the queue data
			nsuccess++
			fmt.Printf("%d,%d,%d\n", total, nsuccess, nfailure)

		case <-ctx.Done():
			// receive abort signal from parent context
			return errors.ToIsyncError(130, "aborted by task")
		}
	}
}
