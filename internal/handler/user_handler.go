package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/jovan/mybanksoal-api/internal/usecase"
	"github.com/jovan/mybanksoal-api/pkg/response"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase}
}

// swagger:model RegisterRequest
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
	APIKey   string `json:"api_key"`
}

// swagger:model LoginRequest
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Register a new user
//
// swagger:route POST /auth/register auth registerUser
//
// Register a new user with username, password, and optional role
//
// Responses:
//   201: successResponse
//   400: errorResponse
func (h *UserHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	if req.Role == "" {
		req.Role = "user"
	}

	user, err := h.userUseCase.Register(req.Username, req.Password, req.Role, req.APIKey)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to register", err.Error())
	}

	return response.Success(c, http.StatusCreated, "User registered successfully", map[string]string{
		"api_key": user.APIKey,
		"username": user.Username,
		"role": user.Role,
	})
}

// Login user
//
// swagger:route POST /auth/login auth loginUser
//
// Login with username and password to get JWT token
//
// Responses:
//   200: successResponse
//   401: errorResponse
func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
	}

	token, err := h.userUseCase.Login(req.Username, req.Password)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
	}

	return response.Success(c, http.StatusOK, "Login successful", map[string]string{
		"token": token,
	})
}
