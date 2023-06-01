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

package modules

import (
	"imshit/aircraftwar/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GetScoreboardRequest struct {
	Mode string `form:"mode" json:"mode" uri:"mode" xml:"mode" binding:"required"`
}

type AddScoreboardRequest struct {
	models.Token
	Score string    `form:"score" json:"score" uri:"score" xml:"score" binding:"required"`
	Mode  string    `form:"mode" json:"mode" uri:"mode" xml:"mode" binding:"required"`
	Time  time.Time `form:"time" json:"time" uri:"time" xml:"time" binding:"required"`
}

type DeleteScoreboardRequest struct {
	models.Token
}

// 获取总排行榜
func GetScoreboard(c *gin.Context) {
	// 检验请求
	request := &GetScoreboardRequest{}
	if err := c.ShouldBind(request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if request.Mode != "easy" && request.Mode != "medium" && request.Mode != "hard" {
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取排行榜
	scores, err := models.GetTopScore(request.Mode, 50)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 格式化
	formatted := []map[string]interface{}{}
	for index, info := range scores {
		if user, err := models.GetUser(info.UserID); err == nil {
			formatted = append(formatted, map[string]interface{}{
				"id":    info.ID,
				"rank":  index,
				"user":  user.Name,
				"score": info.Score,
				"mode":  info.Mode,
				"time":  info.Time,
			})
		}
	}
	// 返回
	c.JSON(http.StatusOK, formatted)
}

// 新增排行榜项目
func AddScoreboard(c *gin.Context) {
	// 检验请求
	request := &AddScoreboardRequest{}
	if err := c.ShouldBind(request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	score, err := strconv.ParseInt(request.Score, 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if score < 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	if request.Mode != "easy" && request.Mode != "medium" && request.Mode != "hard" {
		c.Status(http.StatusBadRequest)
		return
	}
	if request.Time.After(time.Now()) {
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取用户ID
	user_id, err := models.GetUserIDByToken(&request.Token)
	if err != nil {
		log.Printf("Get user id error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 格式化
	scoreItem := &models.Score{
		UserID: user_id,
		Score:  int(score),
		Mode:   request.Mode,
		Time:   request.Time,
	}
	// 增加
	if err := models.SaveScore(scoreItem); err != nil {
		log.Printf("Score saving error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

// 删除排行榜项目
func DeleteScoreboard(c *gin.Context) {
	// 检验请求
	request := &DeleteScoreboardRequest{}
	if err := c.ShouldBind(request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	score_id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if score_id < 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取用户ID
	user_id, err := models.GetUserIDByToken(&request.Token)
	if err != nil {
		log.Printf("Get user id error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取记录
	score, err := models.GetScore(int(score_id))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 核验记录是否是用户产生的
	if score.UserID != user_id {
		c.Status(http.StatusBadRequest)
		return
	}
	// 删除
	if err := models.DeleteScore(score); err != nil {
		log.Printf("Score delete error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}
