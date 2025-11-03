package vserver

import "context"

type ServerOptions = func(server Server) Server

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Register(opts ...ServerOptions) Server
}
