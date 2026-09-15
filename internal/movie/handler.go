package movie

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"streamflix-backend/internal/embedding"
	"github.com/pgvector/pgvector-go"
	"github.com/google/uuid"
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
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "page must be a positive integer",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit must be between 1 and 100",
		})
		return
	}

	movies, total, err := h.service.GetMovies(
		c.Request.Context(),
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
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

func (h *Handler) GetMovie(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid movie id",
		})
		return
	}

	movie, err := h.service.GetMovie(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "movie not found",
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}