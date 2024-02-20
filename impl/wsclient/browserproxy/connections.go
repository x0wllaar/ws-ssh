package browserproxy

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"grisha.xyz/ws-ssh/util"
	"nhooyr.io/websocket"
)

func (d *browserProxyWSDialer) acceptControlConnection(w http.ResponseWriter, r *http.Request) {
	connectionLogger := d.logger.With(slog.String("connection", "control"))
	connectionLogger.Debug("Handling control request", slog.String("remote", r.RemoteAddr), slog.String("url", r.URL.String()))

	websockRaw, err := websocket.Accept(w, r, nil)
	if err != nil {
		connectionLogger.Error("Error accepting websocket connection", util.SlogError(err))
		d.controlAccept <- fmt.Errorf("error accepting websocket connection: %w", err)
		return
	}

	msgType, msgBytes, err := websockRaw.Read(context.Background())
	if err != nil {
		connectionLogger.Error("Error reading client hello", util.SlogError(err))
		d.controlAccept <- fmt.Errorf("error reading client hello: %w", err)
		return
	}
	if msgType != websocket.MessageText {
		connectionLogger.Error("Incorrect client hello message type")
		d.controlAccept <- fmt.Errorf("incorrect client hello message type")
		return
	}
	msgStr := string(msgBytes)
	if msgStr != "BPROXY-HELO" {
		connectionLogger.Error("Incorrect client hello message", util.SlogError(err))
		d.controlAccept <- fmt.Errorf("incorrect client hello message")
		return
	}

	d.controlConn = websockRaw
	d.controlAccept <- nil

	d.logger.Debug("Control connection established")
}

func (d *browserProxyWSDialer) acceptDataConnection(w http.ResponseWriter, r *http.Request) {
	connId := r.PathValue("id")
	connectionLogger := d.logger.With(slog.String("connection", connId))
	connectionLogger.Info("Handling data connection", slog.String("remote", r.RemoteAddr), slog.String("url", r.URL.String()))

	d.wsAcceptConnMapLock.RLock()
	resChan, ok := d.wsAcceptConnMap[connId]
	d.wsAcceptConnMapLock.RUnlock()
	if !ok {
		connectionLogger.Error("Error accepting websocket data connection, id not found in map")
		return
	}

	websockRaw, err := websocket.Accept(w, r, nil)
	if err != nil {
		connectionLogger.Error("Error accepting websocket connection", util.SlogError(err))
		resChan <- acceptedWSConn{nil, fmt.Errorf("error accepting websocket connection: %w", err)}
		return
	}

	resChan <- acceptedWSConn{websockRaw, nil}

}
