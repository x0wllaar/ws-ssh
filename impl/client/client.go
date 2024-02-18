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
	"fmt"
	"io"
	"log/slog"
	"os"

	"grisha.xyz/ws-ssh/util"
	"nhooyr.io/websocket"
)

func connectCmdImpl(logger *slog.Logger, url string, from io.ReadWriter) error {
	localLogger := logger.With(slog.String("to", url))
	localLogger.Info("Starting cleint")

	websockRaw, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		localLogger.Error("Error connecting to websocket", slog.String("error", err.Error()))
		return fmt.Errorf("error connecting to websocket: %w", err)
	}
	defer websockRaw.CloseNow()

	websock := websocket.NetConn(context.Background(), websockRaw, websocket.MessageBinary)
	defer websock.Close()

	err = util.StreamCopy(localLogger, from, websock)
	if err != nil {
		localLogger.Warn("Errors when copying streams", slog.String("error", err.Error()))
	}

	slog.Info("Done copying streams")
	return nil
}

func ConnectCmdImplStdIo(logger *slog.Logger, url string) error {
	logger.Info("Using stdio")
	stdioRw := newConnectedReadWriter(os.Stdin, os.Stdout)
	return connectCmdImpl(logger, url, stdioRw)
}
