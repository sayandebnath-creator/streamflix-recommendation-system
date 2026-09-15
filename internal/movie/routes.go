package movie

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
) {
	router.GET("/movies", handler.GetAllMovies)
	router.GET("/movies/:id", handler.GetMovie)
	router.POST("/movies", handler.CreateMovie)
	router.GET("/search", handler.SearchMovies)
}