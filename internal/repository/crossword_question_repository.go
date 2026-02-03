package repository

import (
	"github.com/jovan/mybanksoal-api/internal/entity"
	"gorm.io/gorm"
)

type CrosswordQuestionRepository interface {
	Create(question *entity.CrosswordQuestion) error
	FindAll(levelID uint) ([]entity.CrosswordQuestion, error)
	FindByID(id uint) (*entity.CrosswordQuestion, error)
	Update(question *entity.CrosswordQuestion) error
	Delete(id uint) error
}

type crosswordQuestionRepository struct {
	db *gorm.DB
}

func NewCrosswordQuestionRepository(db *gorm.DB) CrosswordQuestionRepository {
	return &crosswordQuestionRepository{db}
}

func (r *crosswordQuestionRepository) Create(question *entity.CrosswordQuestion) error {
	if question.QuestionsID != nil {
		var q entity.Question
		if err := r.db.First(&q, *question.QuestionsID).Error; err != nil {
			return err
		}
		question.Clue = q.Content
		question.Answer = q.Answer
	}
	return r.db.Create(question).Error
}

func (r *crosswordQuestionRepository) FindAll(levelID uint) ([]entity.CrosswordQuestion, error) {
	var questions []entity.CrosswordQuestion
	query := r.db.Model(&entity.CrosswordQuestion{})
	if levelID != 0 {
		query = query.Where("level_id = ?", levelID)
	}
	err := query.Find(&questions).Error
	return questions, err
}

func (r *crosswordQuestionRepository) FindByID(id uint) (*entity.CrosswordQuestion, error) {
	var question entity.CrosswordQuestion
	err := r.db.First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *crosswordQuestionRepository) Update(question *entity.CrosswordQuestion) error {
	return r.db.Save(question).Error
}

func (r *crosswordQuestionRepository) Delete(id uint) error {
	return r.db.Delete(&entity.CrosswordQuestion{}, id).Error
}
