package config

import (
	"crawler/baseline/internal/http/controller"
	"crawler/baseline/internal/http/route"
	"crawler/baseline/internal/repository"
	"crawler/baseline/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB     *gorm.DB
	Log    *logrus.Logger
	Config *viper.Viper
}

func Bootstrap(config *BootstrapConfig) *chi.Mux {
	repoRepository := repository.NewRepoRepository(config.Log)
	releaseReposotory := repository.NewReleaseRepository(config.Log)

	repoUsecase := usecase.NewRepoUsecase(config.DB, config.Log, repoRepository)
	releaseUsecase := usecase.NewReleaseUsecase(config.DB, config.Log, releaseReposotory)

	repoController := controller.NewRepoController(config.Log, config.DB, repoUsecase)
	releaseController := controller.NewReleaseController(config.Log, config.DB, releaseUsecase)

	route := route.RouteConfig{
		App:               chi.NewRouter(),
		RepoController:    repoController,
		ReleaseController: releaseController,
	}
	r := route.Setup()
	return r

}
