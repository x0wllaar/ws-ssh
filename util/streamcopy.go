package util

import (
	"fmt"
	"io"
	"log/slog"
)

func StreamCopy(logger *slog.Logger, from io.ReadWriter, to io.ReadWriter) error {
	tfChan := make(chan error)
	ftChan := make(chan error)

	go func() {
		logger.Debug("Starting to -> from copy")
		n, err := io.Copy(from, to)
		if err != nil {
			logger.Debug("To -> from copy ended", slog.Int64("copied", n), slog.String("error", err.Error()))
		} else {
			logger.Debug("To -> from copy ended", slog.Int64("copied", n))
		}
		tfChan <- err
	}()

	go func() {
		logger.Debug("Starting from -> to copy")
		n, err := io.Copy(to, from)
		if err != nil {
			logger.Debug("From -> to copy ended", slog.Int64("copied", n), slog.String("error", err.Error()))
		} else {
			logger.Debug("From -> to copy ended", slog.Int64("copied", n))
		}
		ftChan <- err
	}()

	errTf := <-tfChan
	errFt := <-ftChan

	if (errTf == nil) && (errFt == nil) {
		return nil
	}

	return fmt.Errorf("to to from: %w, from to to: %w", errTf, errFt)
}
