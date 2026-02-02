package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/usecase"
	"github.com/jovan/mybanksoal-api/pkg/response"
)

type QuestionHandler struct {
	questionUseCase usecase.QuestionUseCase
}

func NewQuestionHandler(questionUseCase usecase.QuestionUseCase) *QuestionHandler {
	return &QuestionHandler{questionUseCase}
}

// swagger:model CreateQuestionRequest
type CreateQuestionRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Answer  string `json:"answer"`
}

// swagger:model UpdateQuestionRequest
type UpdateQuestionRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Answer  *string `json:"answer"`
}

// swagger:model UpdateStatusRequest
type UpdateStatusRequest struct {
	Status entity.QuestionStatus `json:"status"`
}

// Create a new question
//
// swagger:route POST /questions questions createQuestion
//
// Create a new question (starts as draft)
//
// Responses:
//   201: successResponse
//   400: errorResponse
func (h *QuestionHandler) Create(c echo.Context) error {
	var req CreateQuestionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	userID := c.Get("user_id").(uint)

	question, err := h.questionUseCase.Create(req.Title, req.Content, req.Answer, userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create question", err.Error())
	}

	// Filter user info if not admin
	role, _ := c.Get("role").(string)
	if role != "admin" {
		question.User = nil
	}

	return response.Success(c, http.StatusCreated, "Question created successfully", question)
}

// Update a question
//
// swagger:route PUT /questions/{id} questions updateQuestion
//
// Update question content
//
// Responses:
//   200: successResponse
//   400: errorResponse
func (h *QuestionHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateQuestionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	question, err := h.questionUseCase.Update(uint(id), req.Title, req.Content, req.Answer)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update question", err.Error())
	}

	// Filter user info if not admin
	role, _ := c.Get("role").(string)
	if role != "admin" {
		question.User = nil
	}

	return response.Success(c, http.StatusOK, "Question updated successfully", question)
}

// Update question status
//
// swagger:route PATCH /questions/{id}/status questions updateQuestionStatus
//
// Update question status (draft, publish, unpublish)
//
// Responses:
//   200: successResponse
//   400: errorResponse
func (h *QuestionHandler) UpdateStatus(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateStatusRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	question, err := h.questionUseCase.UpdateStatus(uint(id), req.Status)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update status", err.Error())
	}

	// Filter user info if not admin
	role, _ := c.Get("role").(string)
	if role != "admin" {
		question.User = nil
	}

	return response.Success(c, http.StatusOK, "Status updated successfully", question)
}

// Delete a question
//
// swagger:route DELETE /questions/{id} questions deleteQuestion
//
// Delete a question by ID
//
// Responses:
//   200: successResponse
//   400: errorResponse
func (h *QuestionHandler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.questionUseCase.Delete(uint(id)); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete question", err.Error())
	}

	return response.Success(c, http.StatusOK, "Question deleted successfully", nil)
}

// Get question by ID
//
// swagger:route GET /questions/{id} questions getQuestion
//
// Get question details
//
// Responses:
//   200: successResponse
//   404: errorResponse
func (h *QuestionHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	question, err := h.questionUseCase.GetByID(uint(id))
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Question not found", err.Error())
	}

	// Filter user info if not admin
	role, _ := c.Get("role").(string)
	if role != "admin" {
		question.User = nil
	}

	return response.Success(c, http.StatusOK, "Question found", question)
}

// Get all questions
//
// swagger:route GET /questions questions listQuestions
//
// Get a list of questions with pagination
//
// Responses:
//   200: successResponse
//   500: errorResponse
func (h *QuestionHandler) GetAll(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}

	questions, err := h.questionUseCase.GetAll(offset, limit)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to fetch questions", err.Error())
	}

	// Filter user info if not admin
	role, _ := c.Get("role").(string)
	if role != "admin" {
		for i := range questions {
			questions[i].User = nil
		}
	}

	return response.Success(c, http.StatusOK, "Questions fetched successfully", questions)
}
