package util

import "log/slog"

func SlogError(err error) slog.Attr {
	return slog.String("error", err.Error())
}
