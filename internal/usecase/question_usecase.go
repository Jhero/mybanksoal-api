package usecase

import (
	"errors"

	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/repository"
)

type QuestionUseCase interface {
	Create(title, content, answer string, createdBy uint) (*entity.Question, error)
	Update(id uint, title, content, answer *string) (*entity.Question, error)
	UpdateStatus(id uint, status entity.QuestionStatus) (*entity.Question, error)
	Delete(id uint) error
	GetByID(id uint) (*entity.Question, error)
	GetAll(offset, limit int) ([]entity.Question, error)
}

type questionUseCase struct {
	questionRepo repository.QuestionRepository
}

func NewQuestionUseCase(questionRepo repository.QuestionRepository) QuestionUseCase {
	return &questionUseCase{questionRepo}
}

func (u *questionUseCase) Create(title, content, answer string, createdBy uint) (*entity.Question, error) {
	question := &entity.Question{
		Title:     title,
		Content:   content,
		Answer:    answer,
		Status:    entity.StatusDraft,
		CreatedBy: createdBy,
	}

	if err := u.questionRepo.Create(question); err != nil {
		return nil, err
	}

	return question, nil
}

func (u *questionUseCase) Update(id uint, title, content, answer *string) (*entity.Question, error) {
	question, err := u.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if title != nil {
		question.Title = *title
	}
	if content != nil {
		question.Content = *content
	}
	if answer != nil {
		question.Answer = *answer
	}

	if err := u.questionRepo.Update(question); err != nil {
		return nil, err
	}

	return question, nil
}

func (u *questionUseCase) UpdateStatus(id uint, status entity.QuestionStatus) (*entity.Question, error) {
	question, err := u.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Validate status transition if needed
	// For now, allow any transition between Draft, Publish, Unpublish
	switch status {
	case entity.StatusDraft, entity.StatusPublish, entity.StatusUnpublish:
		question.Status = status
	default:
		return nil, errors.New("invalid status")
	}

	if err := u.questionRepo.Update(question); err != nil {
		return nil, err
	}

	return question, nil
}

func (u *questionUseCase) Delete(id uint) error {
	return u.questionRepo.Delete(id)
}

func (u *questionUseCase) GetByID(id uint) (*entity.Question, error) {
	return u.questionRepo.FindByID(id)
}

func (u *questionUseCase) GetAll(offset, limit int) ([]entity.Question, error) {
	return u.questionRepo.FindAll(offset, limit)
}
