package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pento-tech-challenge/backend/api/middleware"
	"pento-tech-challenge/backend/timer"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Server holds the dependencies required for the handlers to work.
type Server struct {
	host         string
	port         uint
	logger       *zap.Logger
	db           *sqlx.DB
	router       *http.ServeMux
	timerService timer.Service
}

// Opts bundles the options for creating a Server.
type Opts struct {
	Host         string
	Port         uint
	Logger       *zap.Logger
	DB           *sqlx.DB
	TimerService timer.Service
}

// NewServer creates a new Server based on the provided Opts and registers defined
// routes in the internal router.
func NewServer(opts Opts) *Server {
	srv := &Server{
		host:         opts.Host,
		port:         opts.Port,
		logger:       opts.Logger,
		db:           opts.DB,
		router:       http.NewServeMux(),
		timerService: opts.TimerService,
	}
	srv.setRoutes()
	return srv
}

func (srv *Server) handleSignals(httpServer *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exit := make(chan int)
	go func() {
		for {
			s := <-ch
			switch s {
			case syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT:
				srv.logger.Info(
					"HTTP_GRACEFUL_SHUTDOWN",
					zap.Any("signal", s),
				)
				exit <- 0

				httpServer.Shutdown(context.Background())
			default:
				srv.logger.Info(
					"HTTP_GRACEFUL_SHUTDOWN_FAIL",
					zap.Any("signal", s),
				)
				exit <- 1
			}
		}
	}()

	code := <-exit
	os.Exit(code)
}

// Run starts the underlying HTTP server on the given address. This method will
// block the execution flow until an error is returned.
func (srv *Server) Run() error {
	chain := middleware.NewChain()

	addr := fmt.Sprintf("%s:%d", srv.host, srv.port)

	server := &http.Server{
		Addr:         addr,
		Handler:      chain.Apply(srv.router),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go srv.handleSignals(server)

	return server.ListenAndServe()
}

func (srv *Server) decode(req *http.Request, value interface{}) error {
	// ct := req.Header.Get("Content-Type")
	// if !strings.Contains(ct, "application/json") {
	// 	return errors.New("unsupported media type")
	// }
	defer req.Body.Close()

	return json.NewDecoder(req.Body).Decode(value)
}

func (srv *Server) json(rw http.ResponseWriter, data interface{}, status int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	if data == nil {
		srv.logger.Debug(
			"HTTP_RESPONSE_JSON_EMPTY",
			zap.Int("status", status),
		)
		return
	}

	srv.logger.Debug(
		"HTTP_RESPONSE_JSON",
		zap.Int("status", status),
	)

	enc := json.NewEncoder(rw)
	enc.SetIndent("", "\t")

	if err := enc.Encode(data); err != nil {
		srv.logger.Error(
			"HTTP_RESPONSE_JSON_FAIL",
			zap.Error(err),
		)
		http.Error(rw, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func (srv *Server) jsonErr(rw http.ResponseWriter, message string, status int) {
	type response struct {
		Error string `json:"error"`
	}
	srv.json(rw, response{Error: message}, status)
}
