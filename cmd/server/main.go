// Package main is the entry point for the blog service HTTP server.
// It handles application initialization, dependency injection, and graceful shutdown.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/handler"
	"github.com/artorias501/blog-service/internal/handler/middleware"
	"github.com/artorias501/blog-service/internal/infrastructure/cache"
	"github.com/artorias501/blog-service/internal/infrastructure/database"
	persistence "github.com/artorias501/blog-service/internal/infrastructure/persistence/repository"
	"github.com/artorias501/blog-service/internal/service"
	"github.com/artorias501/blog-service/pkg/config"
	"github.com/artorias501/blog-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Application holds all dependencies for the blog service.
type Application struct {
	Config  *config.Config
	Logger  *slog.Logger
	DB      *gorm.DB
	Redis   *cache.RedisClient
	Router  *gin.Engine
	Server  *http.Server
	Health  *handler.HealthHandler
	Post    *handler.PostHandler
	Tag     *handler.TagHandler
	Comment *handler.CommentHandler
}

func main() {
	// Parse command line flags
	flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Fail fast if required configuration is missing
		fmt.Fprintf(os.Stderr, "ERROR: Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := initializeLogger(cfg)
	logger.SetDefault(log)

	log.Info("starting blog service",
		slog.String("port", cfg.Server.Port),
		slog.String("environment", cfg.Server.Environment),
	)

	// Initialize application
	app, err := InitializeApplication(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.Any("error", err))
		os.Exit(1)
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Info("HTTP server starting", slog.String("address", app.Server.Addr))
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Graceful shutdown: wait for active requests to complete
	if err := app.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("server stopped")
}

// InitializeApplication creates and wires all application dependencies.
// It follows the dependency injection pattern to ensure loose coupling.
func InitializeApplication(cfg *config.Config, log *slog.Logger) (*Application, error) {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := database.NewConnectionWithMigrate(database.Config{
		DSN: cfg.Database.Path,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize Redis client (optional - service can run without Redis)
	redisClient, err := cache.NewRedisClient(cache.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Warn("Redis connection failed, running without cache", slog.Any("error", err))
		// Continue without Redis - service will work but without caching
	}

	// Initialize repositories
	postRepo := persistence.NewPostRepository(db)
	tagRepo := persistence.NewTagRepository(db)
	commentRepo := persistence.NewCommentRepository(db)

	// Initialize cache repositories (nil if Redis is unavailable)
	var postCache repository.PostCacheRepository
	var tagCache repository.TagCacheRepository
	var commentCache repository.CommentCacheRepository
	if redisClient != nil {
		postCache = cache.NewPostCacheRepository(redisClient)
		tagCache = cache.NewTagCacheRepository(redisClient)
		commentCache = cache.NewCommentCacheRepository(redisClient)
	}

	// Initialize services
	postService := service.NewPostService(
		postRepo,
		postCache,
		tagRepo,
		tagCache,
		commentRepo,
		commentCache,
	)
	tagService := service.NewTagService(tagRepo, tagCache)
	commentService := service.NewCommentService(commentRepo, commentCache, postRepo, postCache)

	// Initialize handlers
	var healthRedis *redis.Client
	if redisClient != nil {
		healthRedis = redisClient.Client
	}
	healthHandler := handler.NewHealthHandler(cfg, db, healthRedis)
	postHandler := handler.NewPostHandler(postService)
	tagHandler := handler.NewTagHandler(tagService)
	commentHandler := handler.NewCommentHandler(commentService)

	// Create router with middleware
	router := setupRouter(cfg, log, healthHandler, postHandler, tagHandler, commentHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Application{
		Config:  cfg,
		Logger:  log,
		DB:      db,
		Redis:   redisClient,
		Router:  router,
		Server:  server,
		Health:  healthHandler,
		Post:    postHandler,
		Tag:     tagHandler,
		Comment: commentHandler,
	}, nil
}

// setupRouter configures all HTTP routes and applies middleware chain.
// Routes are organized by resource and documented with comments.
func setupRouter(
	cfg *config.Config,
	log *slog.Logger,
	healthHandler *handler.HealthHandler,
	postHandler *handler.PostHandler,
	tagHandler *handler.TagHandler,
	commentHandler *handler.CommentHandler,
) *gin.Engine {
	router := gin.New()

	// Apply global middleware chain (order matters!)
	// 1. Recovery - must be first to catch panics in all handlers
	router.Use(middleware.Recovery(log))

	// 2. Request logger - logs all incoming requests
	router.Use(middleware.RequestLogger(log))

	// 3. CORS - handles cross-origin requests
	router.Use(middleware.CORS(cfg))

	// 4. Error handler - processes validation and other errors
	router.Use(middleware.ErrorHandler())

	// ============================================================================
	// HEALTH CHECK ROUTES
	// ============================================================================
	// GET /health - Returns service health status
	// Used by load balancers and monitoring systems
	healthGroup := router.Group("/health")
	{
		healthGroup.GET("", healthHandler.Check)
	}

	// ============================================================================
	// API ROUTES (v1)
	// ============================================================================
	// All API routes are prefixed with /api/v1 for versioning
	v1 := router.Group("/api/v1")
	{
		// ========================================================================
		// POST ROUTES
		// ========================================================================
		// GET    /api/v1/posts           - List all posts (paginated)
		// GET    /api/v1/posts/search    - Search posts by keyword
		// POST   /api/v1/posts           - Create a new post (admin)
		// GET    /api/v1/posts/:id       - Get a specific post
		// PUT    /api/v1/posts/:id       - Update a post (admin)
		// DELETE /api/v1/posts/:id       - Delete a post
		// POST   /api/v1/posts/:id/like  - Like a post
		// GET    /api/v1/posts/:id/tags  - Get tags for a post
		// POST   /api/v1/posts/:id/tags  - Add tag to post (admin)
		// DELETE /api/v1/posts/:id/tags/:tag_id - Remove tag from post (admin)
		posts := v1.Group("/posts")
		{
			posts.GET("", postHandler.ListPosts)          // List posts with pagination
			posts.GET("/search", postHandler.SearchPosts) // Search posts by keyword
			posts.GET("/:id", postHandler.GetPost)        // Get post by ID
			posts.DELETE("/:id", postHandler.DeletePost)  // Delete post
			posts.POST("/:id/like", postHandler.LikePost) // Like a post

			// Post tag operations
			posts.GET("/:id/tags", tagHandler.GetTagsByPost) // Get post's tags

			// ====================================================================
			// COMMENT ROUTES (nested under posts)
			// ====================================================================
			// GET  /api/v1/posts/:id/comments       - List comments for a post
			// GET  /api/v1/posts/:id/comments/count - Get comment count for a post
			posts.GET("/:id/comments", commentHandler.ListCommentsByPost)    // List post's comments
			posts.GET("/:id/comments/count", commentHandler.GetCommentCount) // Get comment count
		}
		adminPosts := v1.Group("/posts")
		adminPosts.Use(middleware.AdminAuth(cfg))
		{
			adminPosts.POST("", postHandler.CreatePost)                           // Create new post
			adminPosts.PUT("/:id", postHandler.UpdatePost)                        // Update post
			adminPosts.POST("/:id/tags", postHandler.AddTagToPost)                // Add tag to post
			adminPosts.DELETE("/:id/tags/:tag_id", postHandler.RemoveTagFromPost) // Remove tag from post
		}

		// ========================================================================
		// TAG ROUTES
		// ========================================================================
		// GET    /api/v1/tags         - List all tags (paginated)
		// GET    /api/v1/tags/search  - Search tags by name
		// POST   /api/v1/tags         - Create a new tag
		// GET    /api/v1/tags/:id     - Get a specific tag
		// PUT    /api/v1/tags/:id     - Update a tag
		// DELETE /api/v1/tags/:id     - Delete a tag
		tags := v1.Group("/tags")
		{
			tags.GET("", tagHandler.ListTags)          // List tags with pagination
			tags.GET("/search", tagHandler.SearchTags) // Search tags by name
			tags.POST("", tagHandler.CreateTag)        // Create new tag
			tags.GET("/:id", tagHandler.GetTag)        // Get tag by ID
			tags.PUT("/:id", tagHandler.UpdateTag)     // Update tag
			tags.DELETE("/:id", tagHandler.DeleteTag)  // Delete tag
		}

		// ========================================================================
		// COMMENT ROUTES
		// ========================================================================
		// GET    /api/v1/comments              - List all comments (admin)
		// GET    /api/v1/comments/:id          - Get a specific comment
		// POST   /api/v1/comments              - Create a new comment
		// PUT    /api/v1/comments/:id          - Update a comment
		// DELETE /api/v1/comments/:id          - Delete a comment
		// POST   /api/v1/comments/:id/approve  - Approve a comment (admin)
		// POST   /api/v1/comments/:id/reject   - Reject a comment (admin)
		// POST   /api/v1/comments/:id/spam     - Mark comment as spam (admin)
		// GET    /api/v1/comments/status/:status - List comments by status (admin)
		comments := v1.Group("/comments")
		{
			comments.POST("", commentHandler.CreateComment)       // Create new comment
			comments.GET("/:id", commentHandler.GetComment)       // Get comment by ID
			comments.PUT("/:id", commentHandler.UpdateComment)    // Update comment
			comments.DELETE("/:id", commentHandler.DeleteComment) // Delete comment
		}
		adminComments := v1.Group("/comments")
		adminComments.Use(middleware.AdminAuth(cfg))
		{
			adminComments.GET("", commentHandler.ListComments)                        // List all comments (admin)
			adminComments.GET("/status/:status", commentHandler.ListCommentsByStatus) // List by status (admin)
			adminComments.POST("/:id/approve", commentHandler.ApproveComment)         // Approve comment (admin)
			adminComments.POST("/:id/reject", commentHandler.RejectComment)           // Reject comment (admin)
			adminComments.POST("/:id/spam", commentHandler.MarkCommentAsSpam)         // Mark as spam (admin)
		}
	}

	return router
}

// Shutdown gracefully shuts down the application.
// It waits for active requests to complete before closing connections.
func (a *Application) Shutdown(ctx context.Context) error {
	// Shutdown HTTP server
	if err := a.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	// Close Redis connection
	if a.Redis != nil {
		if err := a.Redis.Close(); err != nil {
			a.Logger.Warn("failed to close Redis connection", slog.Any("error", err))
		}
	}

	// Close database connection
	sqlDB, err := a.DB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			a.Logger.Warn("failed to close database connection", slog.Any("error", err))
		}
	}

	return nil
}

// initializeLogger creates a logger based on configuration.
func initializeLogger(cfg *config.Config) *slog.Logger {
	// Parse log level
	var level slog.Level
	switch cfg.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return logger.NewWithLevel(cfg.Server.Environment, os.Stdout, level)
}
