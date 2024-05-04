package pkg

import (
	"context"
	_ "embed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type Server struct {
	web    *http.Server
	log    *slog.Logger
	exit   chan os.Signal
	state  *atomic.Bool
	banner string
	cron   *Cron
}

func NewServer() *Server {
	server := new(Server)
	server.exit = make(chan os.Signal)
	signal.Notify(server.exit, os.Interrupt)
	return server
}

func (s *Server) LoadWebService(webServer *http.Server) *Server {
	s.web = webServer
	return s
}

func (s *Server) InitLog(sign string) *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})).WithGroup("server").With("info", "ToT", "sign", sign)
	slog.SetDefault(logger)
	s.log = logger
	slog.Info("logging engine")
	return s
}

func (s *Server) InitCron() *Server {
	s.cron = NewCron()
	slog.Info("cron engine")
	return s
}

func (s *Server) Cron() *Cron {
	return s.cron
}

func (s *Server) Start() error {
	defer func() {
		if s.state != nil {
			s.state.Store(false)
		}
	}()
	_ = s.printBanner()
	go func() {
		if s.web != nil {
			slog.Info("ToT Web Service Start!!!")
			_ = s.web.ListenAndServe()
		}
	}()
	state := &atomic.Bool{}
	state.Store(true)

	s.state = state
	<-s.exit
	slog.Info("listen interrupted signal")
	close(s.exit)

	if s.web != nil {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		s.web.RegisterOnShutdown(func() {
			slog.Info("ToT Web Service Shutdown!!!")
			cancelFunc()
		})
		if err := s.web.Shutdown(ctx); err != nil {
			return err
		}
		<-ctx.Done()
	}
	return nil
}

func (s *Server) Kill() {
	s.exit <- os.Interrupt
}

func (s *Server) Wait() {
	for {
		if !s.state.Load() {
			slog.Info("Waiting for...")
			return
		}
		time.After(time.Millisecond * 500)
	}
}

func (s *Server) SetBanner(filePath string) *Server {
	s.banner = filePath
	return s
}

func (s *Server) printBanner() error {
	if s.banner != "" {
		banBytes, err := os.ReadFile(s.banner)
		if err != nil {
			slog.Error("print Banner error", "err", err.Error())
			return err
		}
		_, _ = os.Stdout.Write(banBytes)
	}
	return nil
}
