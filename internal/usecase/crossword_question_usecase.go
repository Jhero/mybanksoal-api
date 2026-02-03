package usecase

import (
	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/repository"
)

type CrosswordQuestionUseCase interface {
	CreateQuestion(question *entity.CrosswordQuestion) error
	GetQuestions(levelID uint) ([]entity.CrosswordQuestion, error)
	GetQuestion(id uint) (*entity.CrosswordQuestion, error)
	UpdateQuestion(id uint, updates map[string]interface{}) (*entity.CrosswordQuestion, error)
	DeleteQuestion(id uint) error
}

type crosswordQuestionUseCase struct {
	repo repository.CrosswordQuestionRepository
}

func NewCrosswordQuestionUseCase(repo repository.CrosswordQuestionRepository) CrosswordQuestionUseCase {
	return &crosswordQuestionUseCase{repo}
}

func (u *crosswordQuestionUseCase) CreateQuestion(question *entity.CrosswordQuestion) error {
	return u.repo.Create(question)
}

func (u *crosswordQuestionUseCase) GetQuestions(levelID uint) ([]entity.CrosswordQuestion, error) {
	return u.repo.FindAll(levelID)
}

func (u *crosswordQuestionUseCase) GetQuestion(id uint) (*entity.CrosswordQuestion, error) {
	return u.repo.FindByID(id)
}

func (u *crosswordQuestionUseCase) UpdateQuestion(id uint, updates map[string]interface{}) (*entity.CrosswordQuestion, error) {
	question, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Simple update logic, in real app might use mapstructure or manual assignment
	if val, ok := updates["number"]; ok {
		question.Number = int(val.(float64))
	}
	if val, ok := updates["clue"]; ok {
		question.Clue = val.(string)
	}
	if val, ok := updates["answer"]; ok {
		question.Answer = val.(string)
	}
	if val, ok := updates["isAcross"]; ok {
		question.IsAcross = val.(bool)
	}
	if val, ok := updates["row"]; ok {
		question.Row = int(val.(float64))
	}
	if val, ok := updates["col"]; ok {
		question.Col = int(val.(float64))
	}
	if val, ok := updates["questions_id"]; ok {
		if val == nil {
			question.QuestionsID = nil
		} else {
			qid := uint(val.(float64))
			question.QuestionsID = &qid
		}
	}
	if val, ok := updates["level_id"]; ok {
		question.LevelID = uint(val.(float64))
	}

	if err := u.repo.Update(question); err != nil {
		return nil, err
	}
	return question, nil
}

func (u *crosswordQuestionUseCase) DeleteQuestion(id uint) error {
	return u.repo.Delete(id)
}
