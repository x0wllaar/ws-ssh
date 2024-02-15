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
	"os"

	"github.com/spf13/cobra"

	"grisha.xyz/ws-ssh/util"
)

var rootCmd = &cobra.Command{
	Use:   "ws-ssh",
	Short: "A small program that forwards TCP connections over websockets",
	Long: `A small program that forwards TCP connections over websockets.
	
Can be very useful if you want to hide your SSH (or other pure TCP) server
behind Cloudflare or other CDN.

To use, on the server:
run
	ws-ssh listen --from 127.0.0.1:8822 --to 127.0.0.1:22
then add to nginx config:
	location /ws-ssh {
    	proxy_pass  http://127.0.0.1:8822/;
	}

On the client:
	ssh -o ProxyCommand="ws-ssh connect --url https://yoursite.com/ws-ssh stdio" yoursite.com
`,
	PreRun: util.LogConfig,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("loglevel", "warn", "Logging severity level")
}
