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

package middleware

import (
	"bytes"
	"imshit/aircraftwar/models"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func AuthMiddleWare() gin.HandlerFunc {
	// 检查用户登录态
	return func(c *gin.Context) {
		token := &models.Token{}
		// 依托gin判断是否读body
		b := binding.Default(c.Request.Method, c.ContentType())
		if bb, ok := b.(binding.BindingBody); ok {
			// 如果读body，则拷贝一份
			body, err := c.GetRawData()
			if err != nil {
				log.Printf("Middleware get raw data error: %v\n", err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			// 读取原先的body（并穷尽）
			if err := bb.BindBody(body, token); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			// 将原先的body还原
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		} else {
			if err := b.Bind(c.Request, token); err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
		if !models.ValidateToken(token) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}
