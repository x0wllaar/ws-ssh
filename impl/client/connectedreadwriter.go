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

import "io"

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
