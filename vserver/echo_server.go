package vserver

import (
	"context"
	"net/http"

	"github.com/flametest/vita/verrors"
	"github.com/flametest/vita/vmiddleware"
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

func NewEchoServer(ctx context.Context, cfg *EchoServerConfig, opts ...ServerOptions) (Server, error) {
	if cfg == nil {
		return nil, verrors.New(500, "cfg is nil")
	}
	e := echo.New()
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(vmiddleware.APILogMiddleware())
	e.Use(vmiddleware.TraceWithRequestId())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "PUT", "POST", "PATCH", "HEAD", "DELETE"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Request-Id", "X-Authorization"},
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
		echoServer = (opt(echoServer)).(*EchoServer)
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

func (e *EchoServer) Register(opts ...ServerOptions) Server {
	for _, opt := range opts {
		*e = *opt(e).(*EchoServer)
	}
	return e
}

func (e *EchoServer) GetEchoServer() *echo.Echo {
	return e.server
}
