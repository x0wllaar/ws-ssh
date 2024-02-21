package browserproxy

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"grisha.xyz/ws-ssh/util"
	"nhooyr.io/websocket"
)

type BpAddrCtx struct{}

type acceptedWSConn struct {
	conn *websocket.Conn
	err  error
}

type browserProxyWSDialer struct {
	logger *slog.Logger

	controlAccept chan error
	controlConn   *websocket.Conn

	wsAcceptConnMapLock sync.RWMutex
	wsAcceptConnMap     map[string](chan acceptedWSConn)

	server *http.Server
}

func NewBrowserProxyWSDialer(logger *slog.Logger) *browserProxyWSDialer {
	return &browserProxyWSDialer{
		logger:              logger,
		controlAccept:       make(chan error),
		controlConn:         nil,
		wsAcceptConnMapLock: sync.RWMutex{},
		wsAcceptConnMap:     map[string]chan acceptedWSConn{},
		server:              nil,
	}
}

func (d *browserProxyWSDialer) Init(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d.logger.Info("Handling root request", slog.String("remote", r.RemoteAddr), slog.String("url", r.URL.String()))
		w.Header().Add("Referrer-Policy", "no-referrer")
		fmt.Fprintln(w, ROOTPAGE_HTML)
	})

	d.controlAccept = make(chan error)
	mux.HandleFunc("/control/", d.acceptControlConnection)
	mux.HandleFunc("/connections/{id}/", d.acceptDataConnection)

	var serverAddr string
	if ctx.Value(BpAddrCtx{}) != nil {
		serverAddr = ctx.Value(BpAddrCtx{}).(string)
	} else {
		serverAddr = "127.0.0.1:8822"
	}

	d.logger.Debug("Starting server")
	server := &http.Server{Handler: mux, Addr: serverAddr}
	serverErrorChan := make(chan error, 2)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			d.logger.Error("Error in local HTTP server", util.SlogError(err))
			serverErrorChan <- fmt.Errorf("error in local HTTP server: %w", err)
		}
		serverErrorChan <- nil
	}()
	d.logger.Info("Please open your browser and go to the listen address to connect", slog.String("listenaddress", fmt.Sprintf("http://%s", serverAddr)))

	select {
	case ctrlConnError := <-d.controlAccept:
		if ctrlConnError != nil {
			d.logger.Error("Error establishing control connection", util.SlogError(ctrlConnError))
			return fmt.Errorf("error establishing control connection: %w", ctrlConnError)
		}
		d.logger.Debug("Control connection established")
		break
	case serverError := <-serverErrorChan:
		if serverError != nil {
			d.logger.Error("Error establishing control connection", util.SlogError(serverError))
			return fmt.Errorf("error establishing control connection: %w", serverError)
		}
		break
	}

	d.server = server
	return nil
}

func (d *browserProxyWSDialer) Close() error {
	return d.server.Shutdown(context.Background())
}

func (d *browserProxyWSDialer) Dial(ctx context.Context, url string) (*websocket.Conn, error) {
	connId := util.StringGuid()
	connectionLogger := d.logger.With(slog.String("dial", connId))

	dataConnChan := make(chan acceptedWSConn)
	d.wsAcceptConnMapLock.Lock()
	d.wsAcceptConnMap[connId] = dataConnChan
	d.wsAcceptConnMapLock.Unlock()

	connectionLogger.Debug("Sending dial request")
	err := d.controlConn.Write(ctx, websocket.MessageText, []byte(fmt.Sprintf("%s||%s", connId, url)))
	if err != nil {
		connectionLogger.Error("Error sending dial request", util.SlogError(err))
		return nil, fmt.Errorf("error sending dial request: %w", err)
	}

	dataConn := <-dataConnChan
	if dataConn.err != nil {
		connectionLogger.Error("Error establishing data connection", util.SlogError(dataConn.err))
		return nil, fmt.Errorf("error establishing data connection: %w", dataConn.err)
	}

	return dataConn.conn, nil
}
