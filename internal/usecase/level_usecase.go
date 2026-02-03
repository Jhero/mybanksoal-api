package usecase

import (
	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/repository"
)

type LevelUseCase interface {
	GetLevels() ([]entity.Level, error)
	GetLevelByID(id uint) (*entity.Level, error)
	CreateLevel(name string) (*entity.Level, error)
	BulkCreateLevels(levels []entity.Level) error
	UpdateLevel(id uint, name string) (*entity.Level, error)
	DeleteLevel(id uint) error
	SubmitLevel(userID, levelID uint, answers map[int]string) (*entity.UserLevelScore, error)
	GetLeaderboard() ([]repository.LeaderboardEntry, error)
}

type levelUseCase struct {
	repo repository.LevelRepository
}

func NewLevelUseCase(repo repository.LevelRepository) LevelUseCase {
	return &levelUseCase{repo}
}

func (u *levelUseCase) GetLevels() ([]entity.Level, error) {
	return u.repo.FindAll()
}

func (u *levelUseCase) GetLevelByID(id uint) (*entity.Level, error) {
	return u.repo.FindByID(id)
}

func (u *levelUseCase) CreateLevel(name string) (*entity.Level, error) {
	level := &entity.Level{
		Name: name,
	}
	if err := u.repo.Create(level); err != nil {
		return nil, err
	}
	return level, nil
}

func (u *levelUseCase) BulkCreateLevels(levels []entity.Level) error {
	for _, level := range levels {
		// We copy the loop variable to avoid pointer issues if we were passing pointers, 
		// but here we pass value to Create which takes pointer.
		// However, to be safe and clean:
		l := level
		if err := u.repo.Create(&l); err != nil {
			return err
		}
	}
	return nil
}

func (u *levelUseCase) UpdateLevel(id uint, name string) (*entity.Level, error) {
	level, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		level.Name = name
	}

	if err := u.repo.Update(level); err != nil {
		return nil, err
	}

	return level, nil
}

func (u *levelUseCase) DeleteLevel(id uint) error {
	// Check if level exists
	_, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(id)
}
