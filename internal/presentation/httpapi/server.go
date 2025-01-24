package httpapi

import (
	"context"
	"fmt"
	"mygo/internal/option"
	"mygo/pkg/logger"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	opt    *option.Options
	server *http.Server
}

func NewServer(opt *option.Options) *Server {
	s := &Server{
		opt: opt,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", opt.Http.Port),
			Handler: newRouter(opt),
		},
	}

	go func() {
		logger.Info(context.Background(), "start server", "listen", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "start server", "error", err)
		}
	}()

	return s
}

func (s *Server) Close(wg *sync.WaitGroup) error {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
