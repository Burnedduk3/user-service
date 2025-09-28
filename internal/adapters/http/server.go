package http

import (
	"context"
	"fmt"
	"user-service/internal/adapters/http/handlers"
	"user-service/internal/adapters/http/middlewares/logging"
	"user-service/internal/adapters/persistence/user_repository"
	"user-service/internal/application/usecases"
	"user-service/internal/config"
	"user-service/internal/infrastructure"
	"user-service/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo        *echo.Echo
	config      *config.Config
	logger      logger.Logger
	connections *infrastructure.DatabaseConnections // Add this
}

func NewServer(cfg *config.Config, log logger.Logger, connections *infrastructure.DatabaseConnections) (*Server, error) {
	e := echo.New()

	// Configure Echo
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		echo:        e,
		config:      cfg,
		logger:      log,
		connections: connections, // Add this
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup routes
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupMiddleware() {
	// Request ID middleware
	s.echo.Use(middleware.RequestID())

	// Replace Echo's logger with our custom Zap logger
	s.echo.Use(logging.ZapLogger(s.logger.With("component", "http")))

	// Recovery middleware
	s.echo.Use(middleware.Recover())

	// Security headers
	s.echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// CORS middleware
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.config.Server.CORS.AllowOrigins,
		AllowMethods: s.config.Server.CORS.AllowMethods,
		AllowHeaders: s.config.Server.CORS.AllowHeaders,
	}))

	// Request timeout middleware
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: s.config.Server.ReadTimeout,
	}))
}

func (s *Server) setupRoutes() {
	// Health check handlers with database connections
	healthHandler := handlers.NewHealthHandler(s.logger, s.connections) // Updated
	userRepo := user_repository.NewGormUserRepository(s.connections.GetGormDB())

	userUseCases := usecases.NewUserUseCases(userRepo, s.logger)

	userHandler := handlers.NewUserHandler(userUseCases, s.logger)
	// API v1 routes
	v1 := s.echo.Group("/api/v1")

	// Health endpoints
	v1.GET("/health", healthHandler.Health)
	v1.GET("/health/ready", healthHandler.Ready)
	v1.GET("/health/live", healthHandler.Live)

	// Metrics endpoint
	v1.GET("/metrics", healthHandler.Metrics)

	users := v1.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.GET("/email/:email", userHandler.GetUserByEmail)
	}
	s.logRegisteredRoutes()
}

func (s *Server) logRegisteredRoutes() {
	s.logger.Info("HTTP routes registered:")
	for _, route := range s.echo.Routes() {
		s.logger.Info("Route registered",
			"method", route.Method,
			"path", route.Path,
			"name", route.Name)
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.logger.Info("Starting HTTP server", "address", address)

	return s.echo.Start(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server...")
	return s.echo.Shutdown(ctx)
}
