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
	"errors"
	"log/slog"

	"github.com/spf13/cobra"
	"grisha.xyz/ws-ssh/util"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a server and forward a stream to it",
	Long: `Connect to a server and forward a stream to it
Currently, only forwarding stdio to websocket is supported.

To use:
	ws-ssh connect stdio --url wss://yoursite.com/ws-ssh
will connect to wss://yoursite.com/ws-ssh and forward stdio to it`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		localLogger := cmd.Context().Value(logger{}).(*slog.Logger)
		localLogger = localLogger.With(slog.String("command", "connect"))
		ctx := context.WithValue(cmd.Context(), logger{}, localLogger)
		cmd.SetContext(ctx)

		urlString, err := cmd.Flags().GetString("url")
		if err != nil {
			localLogger.Error("Error in URL argument", util.SlogError(err))
			return err
		}
		if urlString == "" {
			localLogger.Error("Empty URL argument")
			return errors.New("empty URL argument")
		}

		ctx = context.WithValue(cmd.Context(), urlStr{}, urlString)
		cmd.SetContext(ctx)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.PersistentFlags().StringP("url", "u", "", "The URL connect to")
}
