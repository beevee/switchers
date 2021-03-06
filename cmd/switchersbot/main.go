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
	"github.com/beevee/switchers/firebase"
	"github.com/beevee/switchers/gameprocessor"
	"github.com/beevee/switchers/telegram"
)

var logger log.Logger

func main() {
	var opts struct {
		LogFile       string `short:"l" long:"logfile" description:"log file name (writes to stdout if not specified)" env:"SWITCHERSBOT_LOGFILE"`
		TelegramToken string `short:"t" long:"telegram-token" description:"Telegram token" env:"SWITCHERSBOT_TELEGRAM_TOKEN"`
		FirebaseToken string `short:"f" long:"firebase-token" description:"Firebase token" env:"SWITCHERSBOT_FIREBASE_TOKEN"`
		FirebaseURL   string `short:"u" long:"firebase-url" description:"Firebase URL" env:"SWITCHERSBOT_FIREBASE_URL"`
		TrumpCode     string `short:"d" long:"trump-code" description:"secret command to become Trump" env:"SWITCHERSBOT_TRUMP_CODE"`
		TeamQuorum    int    `short:"q" long:"team-quorum" description:"how many team members must gather to receive an actual task" default:"4"`
		TeamMinSize   int    `short:"s" long:"team-min-size" description:"minimum size of a team (they can be up to 2x-1 bigger, though)" default:"6"`
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

	playerRepository := &firebase.PlayerRepository{
		Repository: firebase.Repository{
			FirebaseToken: opts.FirebaseToken,
			FirebaseURL:   opts.FirebaseURL,
		},
	}

	roundRepository := &firebase.RoundRepository{
		Repository: firebase.Repository{
			FirebaseToken: opts.FirebaseToken,
			FirebaseURL:   opts.FirebaseURL,
		},
	}

	taskRepository := &firebase.TaskRepository{
		Repository: firebase.Repository{
			FirebaseToken: opts.FirebaseToken,
			FirebaseURL:   opts.FirebaseURL,
		},
	}

	bot := &telegram.Bot{
		TelegramToken: opts.TelegramToken,
		Logger:        log.NewContext(logger).With("component", "telegram"),
	}

	gameProcessor := &gameprocessor.GameProcessor{
		TrumpCode:        opts.TrumpCode,
		TeamQuorum:       opts.TeamQuorum,
		TeamMinSize:      opts.TeamMinSize,
		PlayerRepository: playerRepository,
		RoundRepository:  roundRepository,
		TaskRepository:   taskRepository,
		Bot:              bot,
		Logger:           log.NewContext(logger).With("component", "gameprocessor"),
	}

	bot.GameProcessor = gameProcessor

	mustStart(playerRepository)
	mustStart(roundRepository)
	mustStart(taskRepository)
	mustStart(gameProcessor)
	mustStart(bot)

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	logger.Log("msg", "received signal", "signal", <-signalChannel)

	mustStop(bot)
	mustStop(gameProcessor)
	mustStop(playerRepository)
	mustStop(roundRepository)
	mustStop(taskRepository)
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
