package route

import (
	http "crawler/baseline/internal/http/controller"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouteConfig struct {
	App               *chi.Mux
	RepoController    *http.RepoController
	ReleaseController *http.ReleaseController
}

func (c *RouteConfig) Setup() *chi.Mux {
	// c.SetupGuestRoute()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(120 * time.Second))

	r.Route("/api/repos", func(r chi.Router) {
		r.Get("/crawl", c.RepoController.CrawlAllRepos)
		r.Route("/{repoID}", func(r chi.Router) {
			// r.Use(c.RepoController.RepoCtx)
			r.Get("/", c.RepoController.GetRepo)

		})

	})
	r.Route("/api/releases", func(r chi.Router) {
		r.Get("/crawl", c.ReleaseController.CrawlAllReleases)
		r.Route("/{releaseID}", func(r chi.Router) {
			r.Get("/", c.ReleaseController.GetRelease)
		})
	})
	return r
}
