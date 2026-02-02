package repository

import (
	"github.com/jovan/mybanksoal-api/internal/entity"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(question *entity.Question) error
	Update(question *entity.Question) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Question, error)
	FindAll(offset, limit int) ([]entity.Question, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db}
}

func (r *questionRepository) Create(question *entity.Question) error {
	if err := r.db.Create(question).Error; err != nil {
		return err
	}
	// Reload the question with User association
	var reloadedQuestion entity.Question
	if err := r.db.Preload("User").First(&reloadedQuestion, question.ID).Error; err != nil {
		return err
	}
	*question = reloadedQuestion
	return nil
}

func (r *questionRepository) Update(question *entity.Question) error {
	return r.db.Save(question).Error
}

func (r *questionRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Question{}, id).Error
}

func (r *questionRepository) FindByID(id uint) (*entity.Question, error) {
	var question entity.Question
	if err := r.db.Preload("User").First(&question, id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) FindAll(offset, limit int) ([]entity.Question, error) {
	var questions []entity.Question
	if err := r.db.Preload("User").Offset(offset).Limit(limit).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
