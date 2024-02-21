//go:build js && wasm

package main

import (
	"context"
	"log/slog"
	"syscall/js"

	"grisha.xyz/ws-ssh/util"
	"nhooyr.io/websocket"
)

func HandleConnection(uUrl string, pUrl string) {
	localLogger := moduleLogger.With(slog.String("upstream", uUrl), slog.String("downstream", pUrl))

	go func() {
		uWsConn, _, err := websocket.Dial(context.Background(), uUrl, nil)
		if err != nil {
			localLogger.Error("Error connecting to upstream", util.SlogError(err))
			return
		}
		defer uWsConn.CloseNow()
		uConn := websocket.NetConn(context.Background(), uWsConn, websocket.MessageBinary)
		defer uConn.Close()
		localLogger.Info("Upstream connected")

		pWsConn, _, err := websocket.Dial(context.Background(), pUrl, nil)
		if err != nil {
			localLogger.Error("Error connecting to proxy", util.SlogError(err))
			return
		}
		defer pWsConn.CloseNow()
		pConn := websocket.NetConn(context.Background(), pWsConn, websocket.MessageBinary)
		defer pConn.Close()
		localLogger.Info("Proxy connected")

		err = util.StreamCopy(localLogger, pConn, uConn)
		if err != nil {
			localLogger.Warn("Error copying streams", util.SlogError(err))
			return
		}

		localLogger.Info("Done copying streams")
		return
	}()
}

func JSHandleConnection(this js.Value, p []js.Value) interface{} {
	uURL := p[0].String()
	pURL := p[1].String()
	HandleConnection(uURL, pURL)
	return nil
}
