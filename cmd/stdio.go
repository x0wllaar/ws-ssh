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

	"github.com/spf13/cobra"
	"grisha.xyz/ws-ssh/impl/client"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		localLogger := slog.With(slog.String("command", "connect stdio"))
		urlStr, err := cmd.Flags().GetString("url")
		if err != nil {
			localLogger.Error("Error in URL argument", slog.String("error", err.Error()))
			return err
		}
		if urlStr == "" {
			localLogger.Error("Empty URL argument")
			return err
		}
		err = client.ConnectCmdImplStdIo(localLogger, urlStr)
		if err != nil {
			localLogger.Error("Error connecting", slog.String("error", err.Error()))
			return err
		}
		return nil
	},
}

func init() {
	connectCmd.AddCommand(stdioCmd)
}
