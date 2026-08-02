package recommendation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetRecommendations(c *gin.Context) {
	id := c.Param("id")

	limit := 10
	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid limit",
			})
			return
		}
		limit = parsed
	}

	movies, err := h.service.SimilarMovies(
		c.Request.Context(),
		id,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch recommendations",
		})
		return
	}

	c.JSON(http.StatusOK, movies)
}