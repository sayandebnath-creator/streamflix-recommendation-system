package movie

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"streamflix-backend/internal/embedding"
	"github.com/pgvector/pgvector-go"
)

type Handler struct {
	service           *Service
	embeddingService  embedding.Service
}

func NewHandler(
	service *Service,
	embeddingService embedding.Service,
) *Handler {
	return &Handler{
		service:          service,
		embeddingService: embeddingService,
	}
}

func (h *Handler) GetAllMovies(c *gin.Context) {
	movies, err := h.service.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *Handler) CreateMovie(c *gin.Context) {
	var movie Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.CreateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func (h *Handler) SearchMovies(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query is required",
		})
		return
	}

	vector, err := h.embeddingService.GenerateEmbedding(
		c.Request.Context(),
		query,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	movies, err := h.service.SearchSimilarMovies(
		c.Request.Context(),
		pgvector.NewVector(vector),
		10,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movies)
}