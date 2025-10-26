package vserver

import (
	"context"
	"net/http"

	"github.com/flametest/vita/verrors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServerConfig struct {
	Name string `json:"name" yaml:"name"`
	Addr string `json:"addr" yaml:"addr"`
	Env  string `json:"env" yaml:"env"`
}

type EchoServer struct {
	server *echo.Echo
	ctx    context.Context
	cfg    *EchoServerConfig
}

type EchoServerOptions = func(server *EchoServer) *EchoServer

func NewEchoServer(ctx context.Context, cfg *EchoServerConfig, opts ...EchoServerOptions) (Server, error) {
	if cfg == nil {
		return nil, verrors.New(500, "cfg is nil")
	}
	e := echo.New()
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "PUT", "POST", "PATCH", "HEAD", "DELETE"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Request-ID", "X-Authorization"},
		AllowCredentials: true,
	}))
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	echoServer := &EchoServer{
		server: e,
		ctx:    ctx,
		cfg:    cfg,
	}
	for _, opt := range opts {
		echoServer = opt(echoServer)
	}
	return echoServer, nil
}

func (e *EchoServer) Start(ctx context.Context) error {
	if err := e.server.Start(e.cfg.Addr); err != nil {
		return err
	}
	return nil
}

func (e *EchoServer) Shutdown(ctx context.Context) error {
	if err := e.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
