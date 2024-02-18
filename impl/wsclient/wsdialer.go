package wsclient

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"nhooyr.io/websocket"
)

type WSDialer interface {
	Dial(ctx context.Context, url string) (*websocket.Conn, error)
	Init(ctx context.Context) error
	Close() error
}

func GetDialerForSpec(logger *slog.Logger, spec string) (WSDialer, error) {
	spec = strings.ToLower(spec)
	logger.Info("Getting Websocket dialer for spec", slog.String("spec", spec))
	if spec == "normal" {
		return &normalWSDialer{logger: logger.With(slog.String("dialer", "normalDialer"))}, nil
	}
	logger.Error("Could not find Websocket dialer for spec", slog.String("spec", spec))
	return nil, fmt.Errorf("could not find Websocket dialer for spec %s", spec)
}
