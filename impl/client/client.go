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
package client

import (
	"context"
	"log/slog"
	"os"

	"grisha.xyz/ws-ssh/impl/wsclient"
)

func ConnectCmdImplStdIo(logger *slog.Logger, dialer wsclient.WSDialer, url string) error {
	logger.Info("Using stdio")

	logger.Debug("Initializing dialer")
	dialer.Init(context.Background())
	defer dialer.Close()

	stdioRw := newConnectedReadWriter(os.Stdin, os.Stdout)
	return connectCmdImpl(logger, dialer, url, stdioRw)
}
