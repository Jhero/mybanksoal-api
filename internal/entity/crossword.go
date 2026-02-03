package entity

import "time"

type Level struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	Name      string              `json:"name"`
	Questions []CrosswordQuestion `gorm:"foreignKey:LevelID" json:"questions"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type CrosswordQuestion struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LevelID     uint      `json:"level_id"`
	Number      int       `json:"number"`
	Clue        string    `json:"clue"`
	Answer      string    `json:"answer"`
	IsAcross    bool      `json:"isAcross"`
	Row         int       `json:"row"`
	Col         int       `json:"col"`
	QuestionsID *uint     `json:"questions_id"`
	Question    *Question `gorm:"foreignKey:QuestionsID" json:"question,omitempty"`
}

type UserLevelScore struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	LevelID   uint      `json:"level_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserLevelScore) TableName() string {
	return "user_level_scores"
}
