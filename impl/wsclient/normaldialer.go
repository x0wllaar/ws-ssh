package wsclient

import (
	"context"
	"fmt"
	"log/slog"

	"grisha.xyz/ws-ssh/util"
	"nhooyr.io/websocket"
)

type normalWSDialer struct {
	logger *slog.Logger
}

func (d *normalWSDialer) Dial(ctx context.Context, url string) (*websocket.Conn, error) {
	wsConn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		d.logger.Error("Error while dialing websocket", util.SlogError(err))
		return nil, fmt.Errorf("error websocket dialing: %w", err)
	}
	return wsConn, nil
}

func (d *normalWSDialer) Init(ctx context.Context) error {
	return nil
}

func (d *normalWSDialer) Close() error {
	return nil
}
