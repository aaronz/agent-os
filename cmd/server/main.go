package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agent-os/core/internal/config"
	"github.com/agent-os/core/internal/handlers"
	"github.com/agent-os/core/internal/middleware"
	"github.com/agent-os/core/internal/repository"
	"github.com/agent-os/core/internal/services"
	"github.com/agent-os/core/pkg/tracing"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize tracing
	shutdown, err := tracing.Init(cfg.Tracing)
	if err != nil {
		log.Printf("Warning: Failed to initialize tracing: %v", err)
	} else {
		defer shutdown(context.Background())
	}

	// Initialize database
	db, err := repository.NewDatabase(cfg.Database)
	if err != nil {
		log.Printf("Warning: Database not available: %v", err)
		// Continue without database for development
	} else {
		defer db.Close()
	}

	// Initialize repositories
	orgRepo := repository.NewOrganizationRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	intentRepo := repository.NewIntentRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	bidRepo := repository.NewBidRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	memoryRepo := repository.NewMemoryRepository(db)

	// Initialize graph store
	graphStore, err := repository.NewGraphStore(cfg.Graph)
	if err != nil {
		log.Fatalf("Failed to connect to graph store: %v", err)
	}
	defer graphStore.Close()

	// Initialize vector store
	vectorStore, err := repository.NewVectorStore(cfg.VectorDB)
	if err != nil {
		log.Fatalf("Failed to connect to vector store: %v", err)
	}
	defer vectorStore.Close()

	// Initialize event queue
	eventQueue, err := repository.NewEventQueue(cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to connect to event queue: %v", err)
	}
	defer eventQueue.Close()

	// Initialize services
	authService := services.NewAuthService(agentRepo)
	orgService := services.NewOrganizationService(orgRepo, eventQueue)
	agentService := services.NewAgentService(agentRepo, eventQueue)
	intentService := services.NewIntentService(intentRepo, graphStore, eventQueue)
	taskService := services.NewTaskService(taskRepo, agentRepo, bidRepo, artifactRepo, graphStore, eventQueue)
	artifactService := services.NewArtifactService(artifactRepo, eventQueue)
	reviewService := services.NewReviewService(reviewRepo, agentRepo, artifactRepo, taskRepo, eventQueue)
	memoryService := services.NewMemoryService(memoryRepo, vectorStore, eventQueue)
	reputationService := services.NewReputationService(agentRepo, eventQueue)
	// governanceService := services.NewGovernanceService(agentRepo, eventQueue) // TODO: implement

	// Initialize handlers
	orgHandler := handlers.NewOrganizationHandler(orgService)
	agentHandler := handlers.NewAgentHandler(agentService, authService)
	intentHandler := handlers.NewIntentHandler(intentService)
	taskHandler := handlers.NewTaskHandler(taskService)
	artifactHandler := handlers.NewArtifactHandler(artifactService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	memoryHandler := handlers.NewMemoryHandler(memoryService)
	healthHandler := handlers.NewHealthHandler()
	metricsHandler := handlers.NewMetricsHandler()

	// Setup router
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/metrics", metricsHandler.GetMetrics)

	// Rate limiter - 100 requests per minute
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// API routes
	api := http.NewServeMux()

	// Wrap with middleware
	wrapped := middleware.TraceMiddleware(rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add auth middleware
		// For now, skip auth for development
		api.ServeHTTP(w, r)
	})))

	// Organization routes
	api.HandleFunc("POST /api/v1/organizations", orgHandler.Create)
	api.HandleFunc("GET /api/v1/organizations", orgHandler.List)
	api.HandleFunc("GET /api/v1/organizations/{org_id}", orgHandler.Get)
	api.HandleFunc("PUT /api/v1/organizations/{org_id}", orgHandler.Update)
	api.HandleFunc("PUT /api/v1/organizations/{org_id}/governance", orgHandler.UpdateGovernance)
	api.HandleFunc("PUT /api/v1/organizations/{org_id}/status", orgHandler.UpdateStatus)

	// Agent routes
	api.HandleFunc("POST /api/v1/agents", agentHandler.Create)
	api.HandleFunc("GET /api/v1/agents", agentHandler.List)
	api.HandleFunc("GET /api/v1/agents/{agent_id}", agentHandler.Get)
	api.HandleFunc("PUT /api/v1/agents/{agent_id}", agentHandler.Update)
	api.HandleFunc("PUT /api/v1/agents/{agent_id}/status", agentHandler.UpdateStatus)
	api.HandleFunc("GET /api/v1/agents/{agent_id}/reputation", agentHandler.GetReputation)
	api.HandleFunc("GET /api/v1/agents/{agent_id}/activities", agentHandler.GetActivities)

	// Intent routes
	api.HandleFunc("POST /api/v1/intents", intentHandler.Create)
	api.HandleFunc("GET /api/v1/intents", intentHandler.List)
	api.HandleFunc("GET /api/v1/intents/{intent_id}", intentHandler.Get)
	api.HandleFunc("PUT /api/v1/intents/{intent_id}/status", intentHandler.UpdateStatus)
	api.HandleFunc("GET /api/v1/intents/{intent_id}/trace", intentHandler.GetTrace)

	// Planning routes
	api.HandleFunc("GET /api/v1/planning/intent", intentHandler.GetPendingPlanning)
	api.HandleFunc("POST /api/v1/planning/task-graph", intentHandler.SubmitTaskGraph)
	api.HandleFunc("GET /api/v1/planning/task-graph/{graph_id}", intentHandler.GetTaskGraph)
	api.HandleFunc("PUT /api/v1/planning/task-graph/{graph_id}", intentHandler.UpdateTaskGraph)

	// Task routes
	api.HandleFunc("GET /api/v1/tasks", taskHandler.List)
	api.HandleFunc("GET /api/v1/tasks/{task_id}", taskHandler.Get)
	api.HandleFunc("POST /api/v1/tasks/{task_id}/bid", taskHandler.SubmitBid)
	api.HandleFunc("PUT /api/v1/tasks/{task_id}/accept", taskHandler.AcceptTask)
	api.HandleFunc("PUT /api/v1/tasks/{task_id}/progress", taskHandler.ReportProgress)
	api.HandleFunc("POST /api/v1/tasks/{task_id}/fail", taskHandler.ReportFailure)
	api.HandleFunc("GET /api/v1/tasks/{task_id}/artifacts", taskHandler.GetArtifacts)
	api.HandleFunc("GET /api/v1/tasks/{task_id}/reviews", taskHandler.GetReviews)

	// Artifact routes
	api.HandleFunc("POST /api/v1/artifacts", artifactHandler.Create)
	api.HandleFunc("GET /api/v1/artifacts/{artifact_id}", artifactHandler.Get)
	api.HandleFunc("PUT /api/v1/artifacts/{artifact_id}", artifactHandler.Update)
	api.HandleFunc("GET /api/v1/artifacts/{artifact_id}/versions", artifactHandler.ListVersions)
	api.HandleFunc("GET /api/v1/artifacts/{artifact_id}/dependencies", artifactHandler.GetDependencies)

	// Review routes
	api.HandleFunc("GET /api/v1/review/tasks", reviewHandler.GetPendingReviews)
	api.HandleFunc("POST /api/v1/reviews", reviewHandler.Create)
	api.HandleFunc("GET /api/v1/reviews/{review_id}", reviewHandler.Get)

	// Memory routes
	api.HandleFunc("POST /api/v1/memory/search", memoryHandler.Search)
	api.HandleFunc("GET /api/v1/memory", memoryHandler.List)
	api.HandleFunc("GET /api/v1/memory/{memory_id}", memoryHandler.Get)
	api.HandleFunc("PUT /api/v1/memory/{memory_id}", memoryHandler.Update)

	// Human control routes (admin only)
	api.HandleFunc("POST /api/v1/human/intervention", orgHandler.EmergencyIntervention)
	api.HandleFunc("GET /api/v1/audit", handlers.GetAuditTrail)

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      wrapped,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start event processors
	go services.StartEventProcessor(eventQueue, intentService, taskService, memoryService, reputationService)

	// Start background jobs
	go services.StartReputationDecayJob(agentRepo)
	go services.StartTimeoutDetector(intentRepo, taskRepo)

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
