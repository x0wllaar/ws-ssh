//go:build js && wasm

package main

import (
	"log/slog"
	"os"
	"syscall/js"
)

var moduleLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

func main() {

	js.Global().Set("goHandleProxyConnection", js.FuncOf(JSHandleConnection))
	js.Global().Call("onWasmStartedResolve")

	moduleLogger.Info("wasm initialized")

	wait := make(chan bool)
	<-wait
}
