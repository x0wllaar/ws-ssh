package util

import (
	"strings"

	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"

	"log/slog"

	"github.com/spf13/cobra"
)

func LogConfig(cmd *cobra.Command, args []string) {
	lvlStr, err := cmd.Flags().GetString("loglevel")
	if err != nil {
		logrus.Panicf("Error getting log level: %v", err)
	}
	lvlVal := getLevel(lvlStr)
	logrusLogger := logrus.New()
	logger := slog.New(sloglogrus.Option{Level: lvlVal, Logger: logrusLogger}.NewLogrusHandler())
	slog.SetDefault(logger)
}

func getLevel(lvl string) slog.Level {
	lvl = strings.ToLower(lvl)
	switch lvl {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		logrus.Panicf("Incorrect log level: %v", lvl)
	}
	return slog.LevelDebug
}
