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
