package usecase

import (
	"strings"

	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/repository"
)

func (u *levelUseCase) SubmitLevel(userID, levelID uint, answers map[int]string) (*entity.UserLevelScore, error) {
	level, err := u.repo.FindByID(levelID)
	if err != nil {
		return nil, err
	}

	score := 0
	for _, q := range level.Questions {
		userAns, ok := answers[q.Number]
		if ok && strings.EqualFold(strings.TrimSpace(userAns), q.Answer) {
			score += 10 // 10 points per correct answer
		}
	}

	// Check if user already has a score
	existingScore, err := u.repo.FindScore(userID, levelID)
	if err != nil {
		return nil, err
	}

	if existingScore != nil {
		// Update only if new score is higher
		if score > existingScore.Score {
			existingScore.Score = score
			if err := u.repo.SaveScore(existingScore); err != nil {
				return nil, err
			}
		}
		return existingScore, nil
	}

	// Create new score
	newScore := &entity.UserLevelScore{
		UserID:  userID,
		LevelID: levelID,
		Score:   score,
	}

	if err := u.repo.SaveScore(newScore); err != nil {
		return nil, err
	}

	return newScore, nil
}

func (u *levelUseCase) GetLeaderboard() ([]repository.LeaderboardEntry, error) {
	return u.repo.GetLeaderboard()
}
