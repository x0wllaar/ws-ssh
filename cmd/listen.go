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
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"grisha.xyz/ws-ssh/impl/server"
	"grisha.xyz/ws-ssh/util"
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
	PreRun: util.LogConfig,
	Run: func(cmd *cobra.Command, args []string) {
		localLogger := slog.With(slog.String("command", "listen"))
		fromStr, err := cmd.Flags().GetString("from")
		if err != nil {
			localLogger.Error("Error in from argument", slog.String("error", err.Error()))
			os.Exit(1)
		}
		if fromStr == "" {
			localLogger.Error("Empty from argument")
			os.Exit(1)
		}
		toStr, err := cmd.Flags().GetString("to")
		if err != nil {
			localLogger.Error("Error in to argument", slog.String("error", err.Error()))
			os.Exit(1)
		}
		if toStr == "" {
			localLogger.Error("Empty to argument")
			os.Exit(1)
		}
		server.ListenCmdImpl(localLogger, fromStr, toStr)
	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringP("from", "f", "127.0.0.1:8822", "Where to listen for incoming connections")
	listenCmd.Flags().StringP("to", "t", "127.0.0.1:22", "Where to forward the connections")
}