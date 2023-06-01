/**
 * Copyright (c) [2023] [Jim-shop]
 * [AircraftWarBackend] is licensed under Mulan PubL v2.
 * You can use this software according to the terms and conditions of the Mulan PubL v2.
 * You may obtain a copy of Mulan PubL v2 at:
 *          http://license.coscl.org.cn/MulanPubL-2.0
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PubL v2 for more details.
 */

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
