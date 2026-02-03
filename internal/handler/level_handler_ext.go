package handler

import (
	"net/http"
	"strconv"

	"github.com/jovan/mybanksoal-api/pkg/response"
	"github.com/labstack/echo/v4"
)

type SubmitLevelRequest struct {
	Answers map[int]string `json:"answers"`
}

func (h *LevelHandler) SubmitLevel(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid ID", "")
	}

	var req SubmitLevelRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "")
	}

	// Get user ID from context (set by middleware)
	userID := c.Get("user_id").(uint)

	score, err := h.useCase.SubmitLevel(userID, uint(id), req.Answers)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "")
	}

	return response.Success(c, http.StatusOK, "Level submitted successfully", score)
}

func (h *LevelHandler) GetLeaderboard(c echo.Context) error {
	leaderboard, err := h.useCase.GetLeaderboard()
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "")
	}

	return response.Success(c, http.StatusOK, "Leaderboard fetched successfully", leaderboard)
}
