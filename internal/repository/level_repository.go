package repository

import (
	"github.com/jovan/mybanksoal-api/internal/entity"
	"gorm.io/gorm"
)

type LeaderboardEntry struct {
	Username  string `json:"username"`
	TotalScore int    `json:"total_score"`
}

type LevelRepository interface {
	FindAll() ([]entity.Level, error)
	FindByID(id uint) (*entity.Level, error)
	Create(level *entity.Level) error
	Update(level *entity.Level) error
	Delete(id uint) error
	FindScore(userID, levelID uint) (*entity.UserLevelScore, error)
	SaveScore(score *entity.UserLevelScore) error
	GetLeaderboard() ([]LeaderboardEntry, error)
}

type levelRepository struct {
	db *gorm.DB
}

func NewLevelRepository(db *gorm.DB) LevelRepository {
	return &levelRepository{db}
}

func (r *levelRepository) FindAll() ([]entity.Level, error) {
	var levels []entity.Level
	// Preload Questions
	err := r.db.Preload("Questions").Find(&levels).Error
	return levels, err
}

func (r *levelRepository) FindByID(id uint) (*entity.Level, error) {
	var level entity.Level
	err := r.db.Preload("Questions").First(&level, id).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}

func (r *levelRepository) Create(level *entity.Level) error {
	return r.db.Create(level).Error
}

func (r *levelRepository) Update(level *entity.Level) error {
	return r.db.Save(level).Error
}

func (r *levelRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Level{}, id).Error
}

func (r *levelRepository) FindScore(userID, levelID uint) (*entity.UserLevelScore, error) {
	var score entity.UserLevelScore
	err := r.db.Where("user_id = ? AND level_id = ?", userID, levelID).First(&score).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &score, err
}

func (r *levelRepository) SaveScore(score *entity.UserLevelScore) error {
	return r.db.Save(score).Error
}

func (r *levelRepository) GetLeaderboard() ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	err := r.db.Table("user_level_scores").
		Select("users.username, SUM(user_level_scores.score) as total_score").
		Joins("JOIN users ON users.id = user_level_scores.user_id").
		Group("users.username").
		Order("total_score DESC").
		Limit(10).
		Scan(&entries).Error
	return entries, err
}
