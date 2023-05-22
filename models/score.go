package models

import (
	"imshit/aircraftwar/db"
	"time"

	"gorm.io/gorm"
)

type Score struct {
	gorm.Model
	UserID int       `gorm:"not null"`
	Score  int       `gorm:"not null"`
	Mode   string    `gorm:"not null"`
	Time   time.Time `gorm:"not null"`
}

func GetTopScore(mode string, num int) ([]Score, error) {
	scores := []Score{}
	result := db.GetSql().Table("scores").Where("mode = ?", mode).Order("score DESC").Limit(num).Find(&scores)
	return scores, result.Error
}

func SaveScore(score *Score) error {
	result := db.GetSql().Table("scores").Create(score)
	return result.Error
}

func GetScore(id int) (*Score, error) {
	score := &Score{}
	result := db.GetSql().Table("scores").First(score, id)
	return score, result.Error
}

func DeleteScore(score *Score) error {
	result := db.GetSql().Table("scores").Delete(score)
	return result.Error
}
