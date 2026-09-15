package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"streamflix-backend/internal/config"
	"streamflix-backend/internal/database"
	"streamflix-backend/internal/embedding"
	"streamflix-backend/internal/movie"
	"streamflix-backend/internal/database/migration"
	"streamflix-backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"streamflix-backend/internal/recommendation"
	"github.com/gin-contrib/cors"
)

func main() {
	
	// Load environment variables
	cfg := config.Load()

	metrics.Init()

	// Connect to PostgreSQL
	database.Connect(cfg)
	if err := migration.Run(database.DB); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	embeddingService := embedding.NewHTTPService(
		httpClient,
		cfg.EmbeddingServiceURL,
	)

	movieRepository := movie.NewRepository(database.DB)

	movieService := movie.NewService(
		movieRepository,
	)

	movieHandler := movie.NewHandler(
		movieService,
		embeddingService,
	)

	recommendationRepository := recommendation.NewPostgresRepository(
		database.DB,
	)

	recommendationService := recommendation.NewService(
		recommendationRepository,
	)

	recommendationHandler := recommendation.NewHandler(
		recommendationService,
	)

	log.Println("✅ Database migrated")

	// Create Gin router
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://streamflix-recommendation-frontend.vercel.app",},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))

	router.Use(metrics.Middleware())

	router.SetTrustedProxies(nil)

	// Register the routes
	movie.RegisterRoutes(router, movieHandler)

	router.GET(
		"/movies/:id/recommendations",
		recommendationHandler.GetRecommendations,
	)

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"message":   "Streamflix Backend is running 🚀",
			"health":    "/health",
			"movies":    "/movies",
			"search":    "/search?q=...",
			"metrics":   "/metrics",
			"recommend": "/movies/:id/recommendations",
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Streamflix Backend is running 🚀",
		})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}