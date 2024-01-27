package servers

import (
	"context"
	"net/http"
	"time"
)

type HTTPConfigs struct {
	Port     string `yaml:"port"`
	Timeouts struct {
		Read       uint `yaml:"read-ms"`
		Write      uint `yaml:"write-ms"`
		ReadHeader uint `yaml:"read-header-ms"`
		Idle       uint `yaml:"idle-ms"`
	} `yaml:"timeouts"`
}

type HTTPServer struct {
	httpServer *http.Server
}

func NewHTTPServer(cfg *HTTPConfigs, handler http.Handler) *HTTPServer {
	return &HTTPServer{httpServer: &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadTimeout:       time.Duration(cfg.Timeouts.Read) * time.Millisecond,
		ReadHeaderTimeout: time.Duration(cfg.Timeouts.Read) * time.Millisecond,
		WriteTimeout:      time.Duration(cfg.Timeouts.Read) * time.Millisecond,
		IdleTimeout:       time.Duration(cfg.Timeouts.Read) * time.Millisecond,
	}}
}

func (s *HTTPServer) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
