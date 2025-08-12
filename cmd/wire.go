//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"cargo-m/internal/api"
	"cargo-m/internal/repository"
	"cargo-m/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializeApp() *gin.Engine {
	wire.Build(
		repository.NewMavenRepo,

		service.NewMavenService,
		api.NewMavenRepoHandler,
		api.NewRouter,
	)
	return nil
}
