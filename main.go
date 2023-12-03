package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scaffolding/internal/initlize"
	cfg "scaffolding/pkg/config"
	zlog "scaffolding/pkg/log"
	"scaffolding/pkg/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err := initlize.Init(); err != nil {
		panic(err)
	}
	r := gin.New()
	gin.SetMode(cfg.Conf.System.Mode)
	r.Use(gin.Recovery(), middleware.Logger())
	server := &http.Server{
		Addr:    ":" + cfg.Conf.System.Port,
		Handler: r,
	}
	go closeServer(server)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zlog.Panic("failed to start server", zap.Error(err))
	}
}

func closeServer(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	zlog.Info("receive signal", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Conf.System.QuitMaxTime*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Error("failed to shutdown server", zap.Error(err))
	}
	_ = zlog.Logger.Sync()
}
