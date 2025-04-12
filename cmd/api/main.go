// main.go
package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"gobizmanager/internal/app/auth"
	"gobizmanager/internal/app/company"
	"gobizmanager/internal/app/company_user"
	"gobizmanager/internal/app/rbac"
	"gobizmanager/internal/app/role"
	"gobizmanager/internal/app/role/permission"
	"gobizmanager/internal/app/user"
	"gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/migration"
	"gobizmanager/platform/config"
	"gobizmanager/platform/database"
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
	db, err := database.NewDB(cfg.DBPath)
	if err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		return
	}

	// Get underlying SQL DB for migrations
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get SQL DB", zap.Error(err))
		return
	}
	defer sqlDB.Close()

	// Apply migrations
	if err := migration.ApplyMigrations(sqlDB); err != nil {
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
	companyUserRepo := company_user.NewRepository(db, cfg)

	// Initialize handlers
	authHandler := auth.NewHandler(userRepo, jwtManager, msgStore)
	companyHandler := company.NewHandler(companyRepo, rbacRepo, userRepo, roleRepo, permissionRepo, msgStore)
	roleHandler := rbac.NewRoleHandler(rbacRepo, msgStore)
	permissionHandler := rbac.NewPermissionHandler(rbacRepo, msgStore)
	companyUserHandler := company_user.NewHandler(companyUserRepo, rbacRepo, msgStore)
	userHandler := user.NewHandler(userRepo)

	// Create router
	r := chi.NewRouter()

	// Add CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(corsMiddleware.Handler)

	// Add other middleware
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
		r.Mount("/rbac", rbac.Routes(roleHandler, permissionHandler))
		r.Mount("/company-users", company_user.Routes(companyUserHandler))
		r.Mount("/users", user.Routes(userHandler))
	})

	// Start server
	logger.Info("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error("Server failed to start", zap.Error(err))
	}
}
