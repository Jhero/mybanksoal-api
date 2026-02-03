package handler

import (
	"net/http"
	"strconv"

	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/usecase"
	"github.com/jovan/mybanksoal-api/pkg/response"
	"github.com/labstack/echo/v4"
)

type CrosswordQuestionHandler struct {
	useCase usecase.CrosswordQuestionUseCase
}

func NewCrosswordQuestionHandler(u usecase.CrosswordQuestionUseCase) *CrosswordQuestionHandler {
	return &CrosswordQuestionHandler{u}
}

// swagger:model CreateCrosswordQuestionRequest
type CreateCrosswordQuestionRequest struct {
	LevelID     uint   `json:"level_id" validate:"required"`
	Number      int    `json:"number" validate:"required"`
	Clue        string `json:"clue" validate:"required"`
	Answer      string `json:"answer" validate:"required"`
	IsAcross    bool   `json:"isAcross"`
	Row         int    `json:"row" validate:"required"`
	Col         int    `json:"col" validate:"required"`
	QuestionsID *uint  `json:"questions_id"`
}

// swagger:model UpdateCrosswordQuestionRequest
type UpdateCrosswordQuestionRequest struct {
	LevelID     *uint   `json:"level_id,omitempty"`
	Number      *int    `json:"number,omitempty"`
	Clue        *string `json:"clue,omitempty"`
	Answer      *string `json:"answer,omitempty"`
	IsAcross    *bool   `json:"isAcross,omitempty"`
	Row         *int    `json:"row,omitempty"`
	Col         *int    `json:"col,omitempty"`
	QuestionsID *uint   `json:"questions_id,omitempty"`
}

func (h *CrosswordQuestionHandler) Create(c echo.Context) error {
	var req CreateCrosswordQuestionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	question := &entity.CrosswordQuestion{
		LevelID:     req.LevelID,
		Number:      req.Number,
		Clue:        req.Clue,
		Answer:      req.Answer,
		IsAcross:    req.IsAcross,
		Row:         req.Row,
		Col:         req.Col,
		QuestionsID: req.QuestionsID,
	}

	if err := h.useCase.CreateQuestion(question); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to create question", err.Error())
	}

	return response.Success(c, http.StatusCreated, "Crossword question created successfully", question)
}

func (h *CrosswordQuestionHandler) GetAll(c echo.Context) error {
	levelID := 0
	if cid := c.QueryParam("level_id"); cid != "" {
		if id, err := strconv.Atoi(cid); err == nil {
			levelID = id
		}
	}

	questions, err := h.useCase.GetQuestions(uint(levelID))
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to fetch questions", err.Error())
	}

	return response.Success(c, http.StatusOK, "Crossword questions fetched successfully", questions)
}

func (h *CrosswordQuestionHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid ID", "")
	}

	question, err := h.useCase.GetQuestion(uint(id))
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Question not found", err.Error())
	}

	return response.Success(c, http.StatusOK, "Question fetched successfully", question)
}

func (h *CrosswordQuestionHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid ID", "")
	}

	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	question, err := h.useCase.UpdateQuestion(uint(id), req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update question", err.Error())
	}

	return response.Success(c, http.StatusOK, "Question updated successfully", question)
}

func (h *CrosswordQuestionHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid ID", "")
	}

	if err := h.useCase.DeleteQuestion(uint(id)); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete question", err.Error())
	}

	return response.Success(c, http.StatusOK, "Question deleted successfully", nil)
}
