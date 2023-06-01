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

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	User     string `form:"user" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

// 核验身份并签发token
func Login(c *gin.Context) {
	// 检查提交的信息
	request := LoginRequest{}
	if err := c.Bind(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 检查用户是否存在
	user, err := models.QueryUser(request.User)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 检查密码是否正确
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password)); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// 签发新Token
	token, err := models.NewToken(user, c)
	if err != nil {
		log.Printf("Token generate error: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.String(http.StatusOK, token.Token)
}
