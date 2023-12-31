// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"flynoob/bibirt-sock/internal/biz"
	"flynoob/bibirt-sock/internal/conf"
	"flynoob/bibirt-sock/internal/server"
	"flynoob/bibirt-sock/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, logger log.Logger) (*kratos.App, func(), error) {
	authService := service.NewAuthService(confServer)
	clientHandler := biz.NewClientHandler(confServer, authService, logger)
	httpServer := server.NewServer(confServer, clientHandler)
	app := newApp(logger, httpServer)
	return app, func() {
	}, nil
}
