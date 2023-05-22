package modules

import (
	"imshit/aircraftwar/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetScoreboardRequest struct {
	Mode string `form:"mode" json:"mode" uri:"mode" xml:"mode" binding:"required"`
}

type AddScoreboardRequest struct {
	models.Token
	Score int       `form:"score" json:"score" uri:"score" xml:"score" binding:"required"`
	Mode  string    `form:"mode" json:"mode" uri:"mode" xml:"mode" binding:"required"`
	Time  time.Time `form:"time" json:"time" uri:"time" xml:"time" binding:"required"`
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
		formatted = append(formatted, map[string]interface{}{
			"rank":    index,
			"user_id": info.UserID,
			"score":   info.Score,
			"mode":    info.Mode,
			"time":    info.Time,
		})
	}
	// 返回
	c.JSON(http.StatusOK, formatted)
}

// 新增排行榜项目
func AddScoreboard(c *gin.Context) {
	// 检验请求
	request := &AddScoreboardRequest{}
	if err := c.ShouldBind(request); err != nil {
		log.Println("bind")
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	if request.Score < 0 {
		log.Println("score")
		c.Status(http.StatusBadRequest)
		return
	}
	if request.Mode != "easy" && request.Mode != "medium" && request.Mode != "hard" {
		log.Println("mode")
		c.Status(http.StatusBadRequest)
		return
	}
	if request.Time.After(time.Now()) {
		log.Println("time")
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取用户ID
	user_id, err := models.GetUserIDByToken(&request.Token)
	if err != nil {
		log.Println("id")
		c.Status(http.StatusBadRequest)
		return
	}
	// 格式化
	score := &models.Score{
		UserID: user_id,
		Score:  request.Score,
		Mode:   request.Mode,
		Time:   request.Time,
	}
	// 增加
	if err := models.SaveScore(score); err != nil {
		log.Println("incre")
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

// 删除排行榜项目
func DeleteScoreboard(c *gin.Context) {
	//todo
	// id := c.Param("id")
}
