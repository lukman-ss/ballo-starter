package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	ballo "github.com/lukman-ss/ballo"
	"github.com/lukman-ss/ballo-starter/internal/task"
)

// TaskStore defines storage methods needed by TaskHandler.
type TaskStore interface {
	Create(title string) (*task.Task, error)
	GetByID(id string) (*task.Task, error)
	List() ([]*task.Task, error)
	Update(id string, title *string, done *bool) (*task.Task, error)
	Delete(id string) error
}

// TaskHandler handles HTTP endpoints for task management.
type TaskHandler struct {
	store        TaskStore
	maxBodyBytes int64
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(store TaskStore, maxBodyBytes int64) *TaskHandler {
	return &TaskHandler{
		store:        store,
		maxBodyBytes: maxBodyBytes,
	}
}

// ErrorResponse represents standardized JSON error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Create handles POST /api/v1/tasks
func (h *TaskHandler) Create(req *ballo.Request) (*ballo.Response, error) {
	var payload struct {
		Title string `json:"title"`
	}
	decoder := json.NewDecoder(io.LimitReader(req.Request.Body, h.maxBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "request body required"}), nil
		}
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return ballo.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{Error: "request body too large"}), nil
		}
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "malformed JSON or invalid fields"}), nil
	}

	if payload.Title == "" {
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "title is required"}), nil
	}
	if len(payload.Title) > 200 {
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "title must be at most 200 characters"}), nil
	}

	created, err := h.store.Create(payload.Title)
	if err != nil {
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()}), nil
	}

	return ballo.JSON(http.StatusCreated, created), nil
}

// List handles GET /api/v1/tasks
func (h *TaskHandler) List(req *ballo.Request) (*ballo.Response, error) {
	tasks, err := h.store.List()
	if err != nil {
		return ballo.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"}), nil
	}
	if tasks == nil {
		tasks = []*task.Task{}
	}
	return ballo.JSON(http.StatusOK, tasks), nil
}

// GetByID handles GET /api/v1/tasks/:id
func (h *TaskHandler) GetByID(req *ballo.Request) (*ballo.Response, error) {
	id := req.Param("id")
	if id == "" {
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "id required"}), nil
	}

	t, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return ballo.JSON(http.StatusNotFound, ErrorResponse{Error: "task not found"}), nil
		}
		return ballo.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"}), nil
	}

	return ballo.JSON(http.StatusOK, t), nil
}

// Update handles PATCH /api/v1/tasks/:id
func (h *TaskHandler) Update(req *ballo.Request) (*ballo.Response, error) {
	id := req.Param("id")
	if id == "" {
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "id required"}), nil
	}

	var payload struct {
		Title *string `json:"title,omitempty"`
		Done  *bool   `json:"done"`
	}
	if err := json.NewDecoder(io.LimitReader(req.Request.Body, h.maxBodyBytes)).Decode(&payload); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "request body required"}), nil
		}
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return ballo.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{Error: "request body too large"}), nil
		}
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "malformed JSON or invalid fields"}), nil
	}

	if payload.Title == nil && payload.Done == nil {
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "at least one field (title or done) must be provided"}), nil
	}

	if payload.Title != nil {
		trimmed := strings.TrimSpace(*payload.Title)
		if trimmed == "" {
			return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "title cannot be empty"}), nil
		}
		if len(trimmed) > 200 {
			return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "title must be at most 200 characters"}), nil
		}
		payload.Title = &trimmed
	}

	updated, err := h.store.Update(id, payload.Title, payload.Done)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return ballo.JSON(http.StatusNotFound, ErrorResponse{Error: "task not found"}), nil
		}
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()}), nil
	}

	return ballo.JSON(http.StatusOK, updated), nil
}

// Delete handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) Delete(req *ballo.Request) (*ballo.Response, error) {
	id := req.Param("id")
	if id == "" {
		return ballo.JSON(http.StatusBadRequest, ErrorResponse{Error: "id required"}), nil
	}

	err := h.store.Delete(id)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return ballo.JSON(http.StatusNotFound, ErrorResponse{Error: "task not found"}), nil
		}
		return ballo.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error"}), nil
	}

	return ballo.JSON(http.StatusNoContent, nil), nil
}
