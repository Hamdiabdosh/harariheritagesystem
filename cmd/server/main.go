package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qirs-mezgeb/api/internal/auth"
	"github.com/qirs-mezgeb/api/internal/audit"
	"github.com/qirs-mezgeb/api/internal/config"
	"github.com/qirs-mezgeb/api/internal/dashboard"
	"github.com/qirs-mezgeb/api/internal/db"
	"github.com/qirs-mezgeb/api/internal/export"
	"github.com/qirs-mezgeb/api/internal/immovable"
	"github.com/qirs-mezgeb/api/internal/middleware"
	"github.com/qirs-mezgeb/api/internal/movable"
	"github.com/qirs-mezgeb/api/internal/photos"
	"github.com/qirs-mezgeb/api/internal/models"
	"github.com/qirs-mezgeb/api/internal/users"
	"github.com/qirs-mezgeb/api/internal/workflow"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()
	if err := db.RunMigrations(cfg.DBURL); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	pool, err := db.Connect(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	if err := os.MkdirAll(cfg.MediaPath, 0o755); err != nil {
		log.Fatalf("create media directory: %v", err)
	}

	router := setupRouter(cfg, pool)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}

	log.Println("server stopped")
}

func setupRouter(cfg *config.Config, pool *pgxpool.Pool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS(cfg.AllowedOrigins))
	router.Use(middleware.ErrorHandler())

	router.GET("/health", healthHandler(pool))

	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.JWTRefreshSecret)
	authHandler := auth.NewHandler(authService)

	usersRepo := users.NewRepository(pool)
	usersService := users.NewService(usersRepo)
	usersHandler := users.NewHandler(usersService)

	auditRepo := audit.NewRepository(pool)
	auditService := audit.NewService(auditRepo)

	immovableRepo := immovable.NewRepository(pool, auditRepo)
	movableRepo := movable.NewRepository(pool, auditRepo)

	photosRepo := photos.NewRepository(pool)
	photosService := photos.NewService(photosRepo, immovableRepo, movableRepo, cfg.MediaPath)
	photosHandler := photos.NewHandler(photosService)

	immovableService := immovable.NewService(immovableRepo, photosService, auditService)
	immovableHandler := immovable.NewHandler(immovableService)

	movableService := movable.NewService(movableRepo, photosService, auditService)
	movableHandler := movable.NewHandler(movableService)

	workflowRepo := workflow.NewRepository(pool)
	workflowService := workflow.NewService(workflowRepo, auditService, immovableRepo, movableRepo)
	workflowHandler := workflow.NewHandler(workflowService)

	dashboardRepo := dashboard.NewRepository(pool)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	exportService := export.NewService(dashboardService, immovableRepo, movableRepo, photosRepo, cfg.MediaPath)
	exportHandler := export.NewHandler(exportService)

	authParser := middleware.AccessTokenParser(func(token string) (middleware.AuthUser, error) {
		claims, err := authService.ParseAccessToken(token)
		if err != nil {
			return middleware.AuthUser{}, err
		}
		return middleware.AuthUser{
			ID:    claims.UserID,
			Email: claims.Email,
			Role:  claims.Role,
		}, nil
	})

	api := router.Group("/api/v1")
	api.GET("/health", healthHandler(pool))

	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.Refresh)

	authorized := api.Group("")
	authorized.Use(middleware.AuthRequired(authParser))
	authorized.POST("/auth/logout", authHandler.Logout)
	authorized.GET("/users/me", usersHandler.GetMe)
	authorized.PUT("/users/me/language", usersHandler.UpdateMyLanguage)

	managerOnly := authorized.Group("")
	managerOnly.Use(middleware.RequireRole(models.RoleManager))
	managerOnly.GET("/users", usersHandler.List)
	managerOnly.POST("/users", usersHandler.Create)
	managerOnly.PUT("/users/:id", usersHandler.Update)
	managerOnly.DELETE("/users/:id", usersHandler.Delete)

	authorized.GET("/records/immovable", immovableHandler.List)
	authorized.GET("/records/immovable/:id", immovableHandler.GetByID)

	registrarOnly := authorized.Group("")
	registrarOnly.Use(middleware.RequireRole(models.RoleRegistrar))
	registrarOnly.POST("/records/immovable", immovableHandler.Create)
	registrarOnly.PUT("/records/immovable/:id/submit", immovableHandler.Submit)
	registrarOnly.PUT("/records/immovable/:id", immovableHandler.Update)

	authorized.GET("/records/movable", movableHandler.List)
	authorized.GET("/records/movable/:id", movableHandler.GetByID)
	registrarOnly.POST("/records/movable", movableHandler.Create)
	registrarOnly.PUT("/records/movable/:id/submit", movableHandler.Submit)
	registrarOnly.PUT("/records/movable/:id", movableHandler.Update)

	registrarOnly.POST("/records/:type/:id/photos", photosHandler.Upload)
	registrarOnly.DELETE("/records/:type/:id/photos/:photo_id", photosHandler.Delete)

	supervisorOnly := authorized.Group("")
	supervisorOnly.Use(middleware.RequireRole(models.RoleSupervisor))
	supervisorOnly.PUT("/records/:type/:id/review-approve", workflowHandler.ReviewApprove)
	supervisorOnly.PUT("/records/:type/:id/review-return", workflowHandler.ReviewReturn)

	managerWorkflow := authorized.Group("")
	managerWorkflow.Use(middleware.RequireRole(models.RoleManager))
	managerWorkflow.PUT("/records/:type/:id/final-approve", workflowHandler.FinalApprove)
	managerWorkflow.PUT("/records/:type/:id/final-return", workflowHandler.FinalReturn)

	supervisorManager := authorized.Group("")
	supervisorManager.Use(middleware.RequireRole(models.RoleSupervisor, models.RoleManager))
	supervisorManager.POST("/records/:type/:id/comments", workflowHandler.AddComment)
	supervisorManager.GET("/export/records/csv", exportHandler.ExportCSV)

	authorized.GET("/records/:type/:id/comments", workflowHandler.GetComments)
	authorized.GET("/records/:type/:id/history", workflowHandler.GetHistory)
	authorized.GET("/records/:type/:id/pdf", exportHandler.ExportPDF)

	authorized.GET("/dashboard/stats", dashboardHandler.GetStats)
	authorized.GET("/records", dashboardHandler.ListRecords)

	return router
}

func healthHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		status := "ok"
		dbStatus := "connected"
		code := http.StatusOK

		if err := db.Ping(ctx, pool); err != nil {
			dbStatus = "disconnected"
			status = "degraded"
			code = http.StatusServiceUnavailable
		}

		c.JSON(code, gin.H{
			"status": status,
			"db":     dbStatus,
		})
	}
}
