package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/jessevdk/go-flags"

	"github.com/beevee/switchers"
	"github.com/beevee/switchers/telegram"
)

var logger log.Logger

func main() {
	var opts struct {
		LogFile       string `short:"l" long:"logfile" description:"log file name (writes to stdout if not specified)" env:"SWITCHERSBOT_LOGFILE"`
		TelegramToken string `short:"t" long:"telegram-token" description:"Telegram token" env:"SWITCHERSBOT_TELEGRAM_TOKEN"`
		FirebaseToken string `short:"f" long:"firebase-token" description:"Firebase token" env:"SWITCHERSBOT_FIREBASE_TOKEN"`
		FirebaseURL   string `short:"u" long:"firebase-url" description:"Firebase URL" env:"SWITCHERSBOT_FIREBASE_URL"`
	}
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(0)
	}

	if opts.LogFile == "" {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	} else {
		logfile, err := os.OpenFile(opts.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open logfile %s: %s", opts.LogFile, err)
			os.Exit(1)
		}
		defer logfile.Close()
		logger = log.NewLogfmtLogger(log.NewSyncWriter(logfile))
	}
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	logger.Log("msg", "starting program", "pid", os.Getpid())

	bot := &telegram.Bot{
		TelegramToken: opts.TelegramToken,
		Logger:        logger,
	}

	mustStart(bot)

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	logger.Log("msg", "received signal", "signal", <-signalChannel)

	mustStop(bot)
}

func mustStart(service switchers.Service) {
	name := reflect.TypeOf(service)

	logger.Log("msg", "starting service", "name", name)
	if err := service.Start(); err != nil {
		logger.Log("msg", "error starting service", "name", name, "error", err)
		os.Exit(1)
	}
	logger.Log("msg", "started service", "name", name)
}

func mustStop(service switchers.Service) {
	name := reflect.TypeOf(service)

	logger.Log("msg", "stopping service", "name", name)
	if err := service.Stop(); err != nil {
		logger.Log("msg", "error stopping service", "name", name, "error", err)
		os.Exit(1)
	}
	logger.Log("msg", "stopped service", "name", name)
}
