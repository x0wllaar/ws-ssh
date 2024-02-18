/*
Copyright © 2024 Grigorii Khvatskii <gkhvatsk@nd.edu>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"os"
	"strings"

	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"

	"log/slog"

	"github.com/spf13/cobra"
)

func logConfig(cmd *cobra.Command, args []string) {
	lvlStr, err := cmd.Flags().GetString("loglevel")
	if err != nil {
		logrus.Panicf("Error getting log level: %v", err)
	}
	lvlVal := getLevel(lvlStr)
	logrusLogger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	localLogger := slog.New(sloglogrus.Option{Level: lvlVal, Logger: logrusLogger}.NewLogrusHandler())
	ctx := context.WithValue(cmd.Context(), logger{}, localLogger)
	cmd.SetContext(ctx)
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
