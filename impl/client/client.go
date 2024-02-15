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
	"io"
	"log/slog"
	"os"

	"grisha.xyz/ws-ssh/util"
	"nhooyr.io/websocket"
)

type connectedReadWriter struct {
	r io.Reader
	w io.Writer
}

func (rw *connectedReadWriter) Read(p []byte) (n int, err error) {
	return rw.r.Read(p)
}

func (rw *connectedReadWriter) Write(p []byte) (n int, err error) {
	return rw.w.Write(p)
}

func newConnectedReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	rw := connectedReadWriter{r, w}
	return &rw
}

func connectCmdImpl(logger *slog.Logger, url string, from io.ReadWriter) {
	localLogger := logger.With(slog.String("to", url))
	localLogger.Info("Starting cleint")

	websockRaw, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		localLogger.Error("Error connecting to websocket", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer websockRaw.CloseNow()

	websock := websocket.NetConn(context.Background(), websockRaw, websocket.MessageBinary)
	defer websock.Close()

	err = util.StreamCopy(localLogger, from, websock)
	if err != nil {
		localLogger.Warn("Errors when copying streams", slog.String("error", err.Error()))
	}

	slog.Info("Done copying streams")
}

func ConnectCmdImplStdIo(logger *slog.Logger, url string) {
	logger.Info("Using stdio")
	stdioRw := newConnectedReadWriter(os.Stdin, os.Stdout)
	connectCmdImpl(logger, url, stdioRw)
}