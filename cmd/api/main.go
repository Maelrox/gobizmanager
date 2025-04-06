// main.go
package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"gobizmanager/internal/auth"
	"gobizmanager/internal/company"
	"gobizmanager/internal/rbac"
	"gobizmanager/internal/role"
	"gobizmanager/internal/role/permission"
	"gobizmanager/internal/user"
	"gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/migration"
	"gobizmanager/platform/config"
	"gobizmanager/platform/middleware/ratelimit"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize logger first
	if err := logger.InitLogger("bin/logs/app.log"); err != nil {
		panic(err)
	}

	// Initialize database
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		return
	}
	defer db.Close()

	// Apply migrations
	if err := migration.ApplyMigrations(db); err != nil {
		logger.Error("Failed to apply migrations", zap.Error(err))
		return
	}

	// Initialize message store
	msgStore := language.NewMessageStore()

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, 15*time.Minute, 24*time.Hour)

	// Initialize repositories
	userRepo := user.NewRepository(db, cfg)
	rbacRepo := rbac.NewRepository(db)
	companyRepo := company.NewRepository(db, cfg, rbacRepo)
	roleRepo := role.NewRepository(db)
	permissionRepo := permission.NewRepository(db)

	// Initialize handlers
	authHandler := auth.NewHandler(userRepo, jwtManager, msgStore)
	companyHandler := company.NewHandler(companyRepo, rbacRepo, userRepo, roleRepo, permissionRepo, msgStore)
	rbacHandler := rbac.NewHandler(rbacRepo, msgStore)

	// Create router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(context.LanguageMiddleware())
	r.Use(ratelimit.New(cfg.RateLimit))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.RefreshToken)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(jwtManager, msgStore))
		r.Mount("/companies", company.Routes(companyHandler, msgStore))
		r.Mount("/rbac", rbac.Routes(rbacHandler))
	})

	// Start server
	logger.Info("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error("Server failed to start", zap.Error(err))
	}
}
