package handler

import (
	"net/http"
	"strconv"

	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/usecase"
	"github.com/jovan/mybanksoal-api/pkg/response"
	"github.com/labstack/echo/v4"
)

type LevelHandler struct {
	useCase usecase.LevelUseCase
}

func NewLevelHandler(u usecase.LevelUseCase) *LevelHandler {
	return &LevelHandler{u}
}

// GetLevels returns all levels with their crosswords
func (h *LevelHandler) GetLevels(c echo.Context) error {
	levels, err := h.useCase.GetLevels()
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "")
	}
	return response.Success(c, http.StatusOK, "Levels fetched successfully", levels)
}

func (h *LevelHandler) GetLevel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid ID", "")
	}

	level, err := h.useCase.GetLevelByID(uint(id))
	if err != nil {
		return response.Error(c, http.StatusNotFound, "Level not found", "")
	}

	return response.Success(c, http.StatusOK, "Level fetched successfully", level)
}

type CreateLevelRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *LevelHandler) CreateLevel(c echo.Context) error {
	var req CreateLevelRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "")
	}
	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "")
	}

	level, err := h.useCase.CreateLevel(req.Name)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "")
	}

	return response.Success(c, http.StatusCreated, "Level created successfully", level)
}

type UpdateLevelRequest struct {
	Name string `json:"name"`
}

func (h *LevelHandler) UpdateLevel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid ID", "")
	}

	var req UpdateLevelRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "")
	}

	level, err := h.useCase.UpdateLevel(uint(id), req.Name)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "")
	}

	return response.Success(c, http.StatusOK, "Level updated successfully", level)
}

func (h *LevelHandler) DeleteLevel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid ID", "")
	}

	if err := h.useCase.DeleteLevel(uint(id)); err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "")
	}

	return response.Success(c, http.StatusOK, "Level deleted successfully", nil)
}

type ImportQuestion struct {
	Number   int    `json:"number"`
	Clue     string `json:"clue"`
	Answer   string `json:"answer"`
	IsAcross bool   `json:"isAcross"`
	Row      int    `json:"row"`
	Col      int    `json:"col"`
}

type ImportLevelItem struct {
	ID        int              `json:"id"`
	Title     string           `json:"title"`
	Questions []ImportQuestion `json:"questions"`
}

type ImportLevelsRequest struct {
	Levels []ImportLevelItem `json:"levels"`
}

func (h *LevelHandler) ImportLevels(c echo.Context) error {
	var req ImportLevelsRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "")
	}

	var levels []entity.Level
	for _, l := range req.Levels {
		var questions []entity.CrosswordQuestion
		for _, q := range l.Questions {
			questions = append(questions, entity.CrosswordQuestion{
				Number:   q.Number,
				Clue:     q.Clue,
				Answer:   q.Answer,
				IsAcross: q.IsAcross,
				Row:      q.Row,
				Col:      q.Col,
			})
		}

		level := entity.Level{
			Name:      l.Title,
			Questions: questions,
		}
		levels = append(levels, level)
	}

	if err := h.useCase.BulkCreateLevels(levels); err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "")
	}

	return response.Success(c, http.StatusCreated, "Levels imported successfully", nil)
}
