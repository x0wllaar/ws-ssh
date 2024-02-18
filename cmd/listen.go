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
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"grisha.xyz/ws-ssh/impl/server"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for incoming websocket connections and forward them to a TCP port",
	Long: `This command starts a process that listens for incoming websocket connections
and forwards them to a TCP port on this or on another server.

For example:
	ws-ssh listen --from 127.0.0.1:8822 --to 127.0.0.1:22
will listen for incoming websocket connections on http://127.0.0.1:8822 and
forward them to 127.0.0.1:22, enabling ssh connections over websockets`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		localLogger := cmd.Context().Value(logger{}).(*slog.Logger)
		localLogger = localLogger.With(slog.String("command", "listen"))
		ctx := context.WithValue(cmd.Context(), logger{}, localLogger)
		cmd.SetContext(ctx)

		fromString, err := cmd.Flags().GetString("from")
		if err != nil {
			localLogger.Error("Error in from argument", slog.String("error", err.Error()))
			return fmt.Errorf("error in from argument: %w", err)
		}
		if fromString == "" {
			localLogger.Error("Empty from argument")
			return fmt.Errorf("empty from argument")
		}
		ctx = context.WithValue(cmd.Context(), fromStr{}, fromString)
		cmd.SetContext(ctx)

		toString, err := cmd.Flags().GetString("to")
		if err != nil {
			localLogger.Error("Error in to argument", slog.String("error", err.Error()))
			return fmt.Errorf("error in to argument: %w", err)
		}
		if toString == "" {
			localLogger.Error("Empty to argument")
			return fmt.Errorf("empty to argument")
		}
		ctx = context.WithValue(cmd.Context(), toStr{}, toString)
		cmd.SetContext(ctx)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		localLogger := cmd.Context().Value(logger{}).(*slog.Logger)

		fromString := cmd.Context().Value(fromStr{}).(string)
		if fromString == "" {
			panic("empty from string")
		}

		toString := cmd.Context().Value(toStr{}).(string)
		if fromString == "" {
			panic("empty to string")
		}

		err := server.ListenCmdImpl(localLogger, fromString, toString)
		if err != nil {
			localLogger.Error("Error running server", slog.String("error", err.Error()))
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.PersistentFlags().StringP("from", "f", "127.0.0.1:8822", "Where to listen for incoming connections")
	listenCmd.PersistentFlags().StringP("to", "t", "127.0.0.1:22", "Where to forward the connections")
}
