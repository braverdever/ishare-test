package handlers

import (
	"net/http"
	"strconv"

	"ishare-task-api/internal/auth"
	"ishare-task-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	db *gorm.DB
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{
		db: db,
	}
}

// CreateTask creates a new task
// @Summary Create Task
// @Description Creates a new task
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body models.CreateTaskRequest true "Task data"
// @Success 201 {object} models.TaskResponse "Task created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "pending"
	}

	// Create task
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := h.db.Create(task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	// Return task response
	c.JSON(http.StatusCreated, models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

// GetTask retrieves a specific task
// @Summary Get Task
// @Description Retrieves a specific task by ID
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} models.TaskResponse "Task retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task ID is required",
		})
		return
	}

	// Parse UUID
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID format",
		})
		return
	}

	// Get task
	var task models.Task
	if err := h.db.Where("id = ?", taskUUID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve task",
		})
		return
	}

	// Return task response
	c.JSON(http.StatusOK, models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

// UpdateTask updates a specific task
// @Summary Update Task
// @Description Updates a specific task by ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID" format(uuid)
// @Param task body models.UpdateTaskRequest true "Task update data"
// @Success 200 {object} models.TaskResponse "Task updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task ID is required",
		})
		return
	}

	// Parse UUID
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID format",
		})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get existing task
	var task models.Task
	if err := h.db.Where("id = ?", taskUUID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve task",
		})
		return
	}

	// Update task fields
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}

	// Save updated task
	if err := h.db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task",
		})
		return
	}

	// Return updated task response
	c.JSON(http.StatusOK, models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

// DeleteTask deletes a specific task
// @Summary Delete Task
// @Description Deletes a specific task by ID
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} map[string]interface{} "Task deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Task not found"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task ID is required",
		})
		return
	}

	// Parse UUID
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID format",
		})
		return
	}

	// Check if task exists
	var task models.Task
	if err := h.db.Where("id = ?", taskUUID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve task",
		})
		return
	}

	// Delete task
	if err := h.db.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// ListTasks retrieves all tasks with optional filtering and pagination
// @Summary List Tasks
// @Description Retrieves all tasks with optional filtering and pagination
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status" example(pending)
// @Param page query int false "Page number" example(1)
// @Param limit query int false "Items per page" example(10)
// @Success 200 {object} models.TasksResponse "Tasks retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// Get query parameters
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Parse pagination parameters
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Build query
	query := h.db.Model(&models.Task{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count tasks",
		})
		return
	}

	// Get tasks with pagination
	var tasks []models.Task
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve tasks",
		})
		return
	}

	// Convert to response format
	taskResponses := make([]models.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = models.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	// Return response
	c.JSON(http.StatusOK, models.TasksResponse{
		Tasks: taskResponses,
		Total: total,
	})
} 