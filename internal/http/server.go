package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vicodevv/relay/internal/storage"
)

type Server struct {
	router *gin.Engine
	repo   *storage.WorkflowRepository
	port   string
}

func NewServer(repo *storage.WorkflowRepository, port string) *Server {
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	server := &Server{
		router: router,
		repo:   repo,
		port:   port,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		api.GET("/health", s.healthCheck)

		workflows := api.Group("/workflows")
		{
			workflows.POST("", s.createWorkflow)
			workflows.POST("/definitions", s.createDefinition)
			workflows.GET("/:id", s.getWorkflow)
			workflows.GET("", s.listWorkflows)
			workflows.GET("/:id/events", s.getWorkflowEvents)
		}
	}
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	go func() {
		logrus.Infof("Server starting on port %s", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logrus.Info("Server stopped")
	return nil
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}
