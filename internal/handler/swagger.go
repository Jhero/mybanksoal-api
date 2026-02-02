package handler

import "github.com/jovan/mybanksoal-api/pkg/response"

// swagger:parameters registerUser
type RegisterParams struct {
	// Register Request Body
	// in: body
	// required: true
	Body RegisterRequest
}

// swagger:parameters loginUser
type LoginParams struct {
	// Login Request Body
	// in: body
	// required: true
	Body LoginRequest
}

// swagger:parameters createQuestion
type CreateQuestionParams struct {
	// Create Question Request Body
	// in: body
	// required: true
	Body CreateQuestionRequest
}

// swagger:parameters updateQuestion
type UpdateQuestionParams struct {
	// Question ID
	// in: path
	// required: true
	ID int `json:"id"`

	// Update Question Request Body
	// in: body
	// required: true
	Body UpdateQuestionRequest
}

// swagger:parameters updateQuestionStatus
type UpdateStatusParams struct {
	// Question ID
	// in: path
	// required: true
	ID int `json:"id"`

	// Update Status Request Body
	// in: body
	// required: true
	Body UpdateStatusRequest
}

// swagger:parameters getQuestion deleteQuestion
type QuestionIDParams struct {
	// Question ID
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:parameters listQuestions
type ListQuestionsParams struct {
	// Offset
	// in: query
	Offset int `json:"offset"`
	// Limit
	// in: query
	Limit int `json:"limit"`
}

// Success Response
// swagger:response successResponse
type SuccessResponseWrapper struct {
	// in: body
	Body response.Response
}

// Error Response
// swagger:response errorResponse
type ErrorResponseWrapper struct {
	// in: body
	Body response.Response
}
