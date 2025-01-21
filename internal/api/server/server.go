package server

import (
	"context"
	"net/http"
	"nstu/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           addr,
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// RunUntilSignal запускает сервер и ждет сигнала для graceful shutdown
func (s *Server) RunUntilSignal() error {
	// Канал для сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Канал для ошибок
	errChan := make(chan error, 1)

	// Запускаем сервер в горутине
	go func() {
		logger.Log.Info().Msg("Сервер запущен: " + s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Ждем либо сигнал завершения, либо ошибку
	select {
	case <-quit:
		logger.Log.Info().Msg("Получен сигнал завершения")
	case err := <-errChan:
		return err
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Log.Error().Err(err).Msg("Ошибка при остановке сервера")
		return err
	}

	logger.Log.Info().Msg("Сервер успешно остановлен")
	return nil
}
