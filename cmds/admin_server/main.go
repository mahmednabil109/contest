package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/linuxboot/contest/cmds/admin_server/server"
	mongoAdapter "github.com/linuxboot/contest/cmds/admin_server/storage/mongo"
	"github.com/linuxboot/contest/pkg/logging"
	"github.com/linuxboot/contest/pkg/xcontext"
	"github.com/linuxboot/contest/pkg/xcontext/bundles/logrusctx"
	"github.com/linuxboot/contest/pkg/xcontext/logger"
)

var (
	flagSet      *flag.FlagSet
	flagPort     *int
	flagDBURI    *string
	flagLogLevel *string
)

func initFlags(cmd string) {
	flagSet = flag.NewFlagSet(cmd, flag.ContinueOnError)
	flagPort = flagSet.Int("port", 8000, "Port to init the admin server on")
	flagDBURI = flagSet.String("dbURI", "mongodb://localhost:27017", "Database URI")
	flagLogLevel = flagSet.String("logLevel", "debug", "A log level, possible values: debug, info, warning, error, panic, fatal")

}

// exitWithError prints the `err` to the stdErr
// exits with code `code`
func exitWithError(err error, code int) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(code)
}

func main() {
	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs, syscall.SIGINT)

	initFlags(os.Args[0])
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		exitWithError(err, 1)
	}

	logLevel, err := logger.ParseLogLevel(*flagLogLevel)
	if err != nil {
		exitWithError(err, 1)
	}

	storage, err := mongoAdapter.NewMongoStorage(*flagDBURI)
	if err != nil {
		exitWithError(err, 1)
	}

	ctx, cancel := xcontext.WithCancel(logrusctx.NewContext(logLevel, logging.DefaultOptions()...))
	defer cancel()

	go func() {
		<-sigs
		cancel()
	}()

	if err := server.Serve(ctx, *flagPort, storage); err != nil {
		exitWithError(fmt.Errorf("server err: %w", err), 1)
	}
}
