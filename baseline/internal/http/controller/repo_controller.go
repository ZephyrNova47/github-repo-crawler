package controller

import (
	"context"
	"crawler/baseline/internal/entity"
	"crawler/baseline/internal/model"
	"crawler/baseline/internal/repository"
	"crawler/baseline/internal/scrape"
	"crawler/baseline/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepoController struct {
	log         *logrus.Logger
	db          *gorm.DB
	repoUsecase *usecase.RepoUsecase
}

func NewRepoController(log *logrus.Logger, db *gorm.DB,
	repoUsecase *usecase.RepoUsecase) *RepoController {
	return &RepoController{
		log:         log,
		db:          db,
		repoUsecase: repoUsecase,
	}
}

func (c *RepoController) RepoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoID := chi.URLParam(r, "repoID")
		repoEntity := &entity.Repository{}
		repoRepository := repository.NewRepoRepository(c.log)
		err := repoRepository.FindById(c.db, repoEntity, repoID)
		if err != nil {
			log.Error("Error finding repo: ", err)
			http.Error(w, "Repo not found", http.StatusNotFound)
			return
		}
		repoResponse := model.RepoResponse{
			ID:       repoEntity.ID,
			RepoName: repoEntity.RepoName,
			UserName: repoEntity.UserName,
		}
		ctx := context.WithValue(r.Context(), "repo", repoResponse)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *RepoController) GetRepo(w http.ResponseWriter, r *http.Request) {
	// Extract repoID from URL parameters
	repoID, _ := strconv.Atoi(chi.URLParam(r, "repoID"))

	c.log.Infof("Fetching repository with ID: %d", repoID)

	// Create repository instance
	repoRepository := repository.NewRepoRepository(c.log)

	// Find repository by ID
	repoEntity := &entity.Repository{}
	err := repoRepository.FindById(c.db, repoEntity, repoID)

	if err != nil {
		c.log.WithError(err).Errorf("Error finding repository with ID %d", repoID)
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	// Convert entity to response model
	repoResponse := &model.RepoResponse{
		ID:       repoEntity.ID,
		RepoName: repoEntity.RepoName,
		UserName: repoEntity.UserName,
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repoResponse); err != nil {
		c.log.WithError(err).Error("Error encoding repository response")
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}
}

func (c *RepoController) CrawlAllRepos(w http.ResponseWriter, r *http.Request) {
	// Ví dụ: tạo sẵn danh sách repo cần crawl
	fmt.Println("Starting to scrape top repositories from gitstar-ranking.com")
	repos, err := scrape.CrawlAllRepos()
	if err != nil {
		c.log.WithError(err).Error("Error crawling repos")
		http.Error(w, "Failed to crawl repos", http.StatusInternalServerError)
		return
	}

	responseData := make([]*model.RepoResponse, 0, len(repos))
	for _, repo := range repos {
		repoResponse, err := c.repoUsecase.Create(r.Context(), repo)
		if err != nil {
			c.log.WithError(err).Error("Error creating repo")
			continue
		}
		c.log.Infof("Created repo: %s", repoResponse.RepoName)
		responseData = append(responseData, repoResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(model.WebResponse[[]*model.RepoResponse]{
		Data: responseData,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// func (c *RepoController) CrawlAll(ctx *fiber.Ctx) error {
// 	limit := 5000
// 	repos := make([]*model.CreateRepoRequest, 0, limit)
// 	reponsesData := make([]*model.RepoResponse, 0, limit)
// 	for i := 0; i < limit; i++ {
// 		repo := repos[i]

// 		repoUsecase := usecase.NewRepoUsecase(
// 			c.db, c.log, c.repoRepo,
// 		)
// 		repoResponse, err := repoUsecase.Create(
// 			ctx.UserContext(), repo)
// 		reponsesData = append(reponsesData, repoResponse)
// 		if err != nil {
// 			log.Error("Error creating repo: ", err)
// 			continue
// 		}

// 		log.Infof("Created repo: %s", repoResponse.Name)
// 	}
// 	return ctx.JSON(model.WebResponse[[]*model.RepoResponse]{Data: reponsesData})
// }
