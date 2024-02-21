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
	"log/slog"

	"github.com/spf13/cobra"
	"grisha.xyz/ws-ssh/impl/client"
	"grisha.xyz/ws-ssh/impl/wsclient"
	"grisha.xyz/ws-ssh/util"
)

// stdioCmd represents the stdio command
var stdioCmd = &cobra.Command{
	Use:   "stdio",
	Short: "Forward stdio to a websocket connection",
	Long: `Connects to a websocket server, then copies bytes from stdin to
the connection and from the connection to stdout.

This can be useful for the ProxyCommand option in SSH.

To use:
	ws-ssh connect stdio --url https://yoursite.com/ws-ssh
To use with SSH:
	ssh -o ProxyCommand="ws-ssh connect --url https://yoursite.com/ws-ssh stdio" yoursite.com
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		localLogger := cmd.Context().Value(logger{}).(*slog.Logger)
		localLogger = localLogger.With(slog.String("command", "connect stdio"))
		ctx := context.WithValue(cmd.Context(), logger{}, localLogger)
		cmd.SetContext(ctx)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		localLogger := cmd.Context().Value(logger{}).(*slog.Logger)

		urlString := cmd.Context().Value(urlStr{}).(string)
		if urlString == "" {
			panic("empty URL string")
		}

		wsDialer := cmd.Context().Value(wsDialerImpl{}).(wsclient.WSDialer)
		if wsDialer == nil {
			panic("nil WS dialer")
		}

		err := client.ConnectCmdImplStdIo(cmd.Context(), localLogger, wsDialer, urlString)
		if err != nil {
			localLogger.Error("Error connecting", util.SlogError(err))
			return err
		}
		return nil
	},
}

func init() {
	connectCmd.AddCommand(stdioCmd)
}
