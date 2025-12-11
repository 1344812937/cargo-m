//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"cargo-m/internal/api"
	"cargo-m/internal/config"
	"cargo-m/internal/core"
	"cargo-m/internal/proxy"
	"cargo-m/internal/repository"
	"cargo-m/internal/service"
	"cargo-m/internal/tasks"

	"github.com/google/wire"
)

func InitializeApp() *core.Application {
	wire.Build(
		config.LoadApplicationConfig,

		repository.NewDataSource,
		repository.NewMavenRepo,

		service.NewMavenService,

		api.NewMavenRepoHandler,
		api.NewBlueCat,
		api.NewRouter,

		tasks.NewCronTask,
		proxy.NewSocksProxy,

		core.NewApplication,
	)
	return nil
}
