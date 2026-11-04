package app

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	ballo "github.com/lukman-ss/ballo"
	ballorender "github.com/lukman-ss/ballo-render"
	"github.com/lukman-ss/ballo-starter/internal/config"
	"github.com/lukman-ss/ballo-starter/internal/http/handler"
	"github.com/lukman-ss/ballo-starter/internal/task"
)

// App encapsulates the Ballo engine, configuration, and app components.
type App struct {
	cfg    *config.Config
	engine *ballo.Engine
	ready  *atomic.Bool
	logger *slog.Logger
}

// New constructs and wires the Ballo application.
func New(cfg *config.Config) (*App, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logger := newLogger(cfg.LogLevel)

	ready := &atomic.Bool{}
	ready.Store(false)

	engine := ballo.New()
	// Initialize renderer for HTML templates
	renderer := ballorender.New(ballorender.Config{Directory: "templates"})
	engine.SetRenderer(renderer)

	// Global middleware
	engine.Use(ballo.LoggerMiddleware(), ballo.RecoveryMiddleware())

	// Custom 404 handler
	engine.NotFound(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusNotFound, map[string]string{
			"error": "resource not found",
		}), nil
	})

	// Custom 405 handler
	engine.MethodNotAllowed(func(req *ballo.Request) (*ballo.Response, error) {
		return ballo.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		}), nil
	})

	// Register health handlers
	healthHandler := handler.NewHealthHandler(ready)
	engine.GET("/healthz", healthHandler.Healthz)
	engine.GET("/readyz", healthHandler.Readyz)

	// Register task handlers
	taskStore := task.NewStore()
	taskHandler := handler.NewTaskHandler(taskStore, cfg.MaxBodyBytes)

	apiV1 := engine.Group("/api/v1")
	apiV1.POST("/tasks", taskHandler.Create)
	apiV1.GET("/tasks", taskHandler.List)
	apiV1.GET("/tasks/:id", taskHandler.GetByID)
	apiV1.PATCH("/tasks/:id", taskHandler.Update)
	apiV1.DELETE("/tasks/:id", taskHandler.Delete)

	// Example HTML rendering via ballo-render
	apiV1.GET("/welcome", func(req *ballo.Request) (*ballo.Response, error) {
		r := engine.Render()
		if r == nil {
			return ballo.JSON(http.StatusInternalServerError, map[string]string{"error": "renderer not configured"}), nil
		}
		var b strings.Builder
		if err := r.Render(&b, "welcome.html", map[string]string{"Title": "Hello Ballo Starter"}); err != nil {
			return nil, err
		}
		return ballo.HTML(http.StatusOK, b.String()), nil
	})

	return &App{
		cfg:    cfg,
		engine: engine,
		ready:  ready,
		logger: logger,
	}, nil
}

// Engine returns the underlying Ballo engine (implements http.Handler).
func (a *App) Engine() *ballo.Engine {
	return a.engine
}

// SetReady sets the application readiness state.
func (a *App) SetReady(ready bool) {
	a.ready.Store(ready)
}

// Logger returns the application structured logger.
func (a *App) Logger() *slog.Logger {
	return a.logger
}

func newLogger(levelStr string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}
