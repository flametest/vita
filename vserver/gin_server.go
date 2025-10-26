package vserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinServerConfig struct {
	Name string `json:"name" yaml:"name"`
	Addr string `json:"addr" yaml:"addr"`
}

type GinServer struct {
	server *http.Server
	cfg    *GinServerConfig
}

type GinServerOptions = func(server *GinServer) *GinServer

func NewGinServer(ctx context.Context, cfg GinServerConfig, opts ...GinServerOptions) (Server, error) {
	r := gin.Default()
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}
	ginServer := &GinServer{
		server: srv,
	}
	for _, opt := range opts {
		ginServer = opt(ginServer)
	}
	return ginServer, nil
}

func (g *GinServer) Start(ctx context.Context) error {
	if err := g.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (g *GinServer) Shutdown(ctx context.Context) error {
	if err := g.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
