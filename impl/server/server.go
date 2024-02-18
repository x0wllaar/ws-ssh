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
package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"grisha.xyz/ws-ssh/util"
	"nhooyr.io/websocket"
)

func ListenCmdImpl(logger *slog.Logger, from string, to string) error {
	localLogger := logger.With(slog.String("from", from), slog.String("to", to))
	localLogger.Info("Starting listener")

	handlerFn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connectionLogger := localLogger.With(slog.String("connection", util.StringGuid()))
		connectionLogger.Info("Handling connection", slog.String("remote", r.RemoteAddr), slog.String("url", r.URL.String()))

		websockRaw, err := websocket.Accept(w, r, nil)
		if err != nil {
			connectionLogger.Error("Error accepting websocket connection", util.SlogError(err))
			return
		}
		defer websockRaw.CloseNow()
		websock := websocket.NetConn(context.Background(), websockRaw, websocket.MessageBinary)
		defer websock.Close()

		tosock, err := net.Dial("tcp", to)
		if err != nil {
			connectionLogger.Error("Error connecting to downstream server", util.SlogError(err))
			websockRaw.Close(websocket.StatusInternalError, "Error connecting to downstream server")
			return
		}
		defer tosock.Close()

		err = util.StreamCopy(connectionLogger, websock, tosock)
		if err != nil {
			connectionLogger.Warn("Errors when copying streams", util.SlogError(err))
		}

		slog.Info("Done copying streams")

	})

	err := http.ListenAndServe(from, handlerFn)
	if err != nil {
		localLogger.Error("Error in HTTP server", util.SlogError(err))
		return err
	}

	return nil
}
