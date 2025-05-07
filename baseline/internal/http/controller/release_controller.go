package controller

import (
	"crawler/baseline/internal/entity"
	"crawler/baseline/internal/model"
	"crawler/baseline/internal/repository"
	"crawler/baseline/internal/scrape"
	"crawler/baseline/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReleaseController struct {
	log            *logrus.Logger
	db             *gorm.DB
	releaseUsecase *usecase.ReleaseUsecase
}

func NewReleaseController(log *logrus.Logger, db *gorm.DB,
	releaseUsecase *usecase.ReleaseUsecase) *ReleaseController {
	return &ReleaseController{
		log:            log,
		db:             db,
		releaseUsecase: releaseUsecase,
	}
}

func (c *ReleaseController) GetRelease(w http.ResponseWriter, r *http.Request) {
	releaseID, _ := strconv.Atoi(chi.URLParam(r, "releaseID"))

	c.log.Infof("Fetching releasesitory with ID: %d", releaseID)

	releasereleasesitory := repository.NewReleaseRepository(c.log)

	releaseEntity := &entity.Release{}
	err := releasereleasesitory.FindById(c.db, releaseEntity, releaseID)

	if err != nil {
		c.log.WithError(err).Errorf("Error finding releasesitory with ID %d", releaseID)
		http.Error(w, "releasesitory not found", http.StatusNotFound)
		return
	}

	releaseResponse := &model.ReleaseResponse{
		ID:      releaseEntity.ID,
		TagName: releaseEntity.TagName,
		Content: releaseEntity.Content,
		RepoID:  releaseEntity.RepoID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(releaseResponse); err != nil {
		c.log.WithError(err).Error("Error encoding releasesitory response")
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}
}

func (c *ReleaseController) CrawlAllReleases(w http.ResponseWriter, r *http.Request) {
	repoEntities := []entity.Repository{}
	repoRepository := repository.NewRepoRepository(c.log)
	err := repoRepository.FindAll(c.db, &repoEntities)
	if err != nil {
		c.log.WithError(err).Error("Error fetching all releases")
		http.Error(w, "Error fetching releases", http.StatusInternalServerError)
		return
	}
	releaseReponses := make([]*model.ReleaseResponse, 0, len(repoEntities))
	for _, repo := range repoEntities {
		repoOwner := repo.UserName
		repoName := repo.RepoName
		repoID := repo.ID
		contents := scrape.CrawlReleases(repoOwner, repoName)
		for tag, content := range contents {
			releaseRequest := &model.CreateReleaseRequest{
				TagName: tag,
				Content: content,
				RepoID:  repoID,
			}
			releaseResponse, err := c.releaseUsecase.Create(
				r.Context(), releaseRequest)
			if err != nil {
				c.log.Error("Error creating release: ", err)
				continue
			}
			c.log.Infof("Created release: %s", releaseResponse.TagName)
			releaseReponses = append(releaseReponses, releaseResponse)
		}

	}

	c.log.Info("Crawling all releases completed.")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(model.WebResponse[[]*model.ReleaseResponse]{
		Data: releaseReponses,
	}); err != nil {
		c.log.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// func (c *releaseController) CrawlAll(ctx *fiber.Ctx) error {
// 	limit := 5000
// 	releases := make([]*model.CreatereleaseRequest, 0, limit)
// 	releasensesData := make([]*model.RepoResponse, 0, limit)
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
