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

package daemon

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type acquireRoomRequest struct {
	master *PairingDaemon
	h0     *pairingHandler
	h1     *pairingHandler
}

type acquireRoomResponce struct {
	request *acquireRoomRequest
	roomID  int
	err     error
}

type FightingDaemon struct {
	rooms   map[int]*Room
	request chan *acquireRoomRequest
	connect chan *fightingHandler
}

var fightingDaemon *FightingDaemon = nil

func GetFightingDaemon() *FightingDaemon {
	if fightingDaemon == nil {
		fightingDaemon = &FightingDaemon{
			rooms:   make(map[int]*Room),
			request: make(chan *acquireRoomRequest, 16),
			connect: make(chan *fightingHandler, 16),
		}
	}
	return fightingDaemon
}

// 获取房间
func (d *FightingDaemon) acquire(request *acquireRoomRequest) {
	for i := 0; i < viper.GetInt("socket.maxRoomNum"); i++ {
		_, ok := d.rooms[i]
		if !ok {
			handler := make(map[uint]*fightingHandler)
			handler[request.h0.user.ID] = nil
			handler[request.h1.user.ID] = nil
			room := &Room{
				mode:         request.h0.mode,
				isStart:      false,
				handler:      handler,
				activateTime: time.Now(),
				close:        make(chan bool, 16),

				connect: make(chan *fightingHandler, 16),
				upload:  make(chan *uploadMessage, 16),
				leave:   make(chan *fightingHandler, 16),
			}
			d.rooms[i] = room
			request.master.room <- &acquireRoomResponce{
				request: request,
				roomID:  i,
				err:     nil,
			}
			go room.Run()
			return
		}
	}
	request.master.room <- &acquireRoomResponce{
		request: request,
		roomID:  0,
		err:     errors.New("exceed max room num"),
	}
}

// 连接房间
func (d *FightingDaemon) dispatch(handler *fightingHandler) {
	room, ok := d.rooms[handler.roomId]
	if !ok {
		bytes, _ := json.Marshal(map[string]any{
			"type": "err",
			"msg":  "no room",
		})
		handler.msg <- bytes
		return
	}
	room.connect <- handler
}

// 清理无用房间
func (d *FightingDaemon) clean() {
	for id, room := range d.rooms {
		if room.activateTime.Add(viper.GetDuration("socket.cleanRoomInterval")).Before(time.Now()) {
			delete(d.rooms, id)
			room.close <- true
			continue
		}
		noConnect := true
		for _, handler := range room.handler {
			if handler != nil {
				noConnect = false
				break
			}
		}
		if noConnect {
			delete(d.rooms, id)
			room.close <- true
		}
	}
}

// 主循环
func (d *FightingDaemon) Run() {
	ticker := time.NewTicker(viper.GetDuration("socket.cleanRoomInterval"))
	defer ticker.Stop()
	for {
		select {
		case request := <-d.request:
			d.acquire(request)
		case handler := <-d.connect:
			d.dispatch(handler)
		case <-ticker.C:
			d.clean()
		}
	}
}

// 用户接口
func (d *FightingDaemon) Connect(ws *websocket.Conn, userId int, roomId int) {
	handler := &fightingHandler{
		ws:     ws,
		userId: userId,
		roomId: roomId,
		room:   nil,
		msg:    make(chan []byte, 16),
	}
	go handler.Write()
	go handler.Read()
	fightingDaemon.connect <- handler
}

type uploadMessage struct {
	from *fightingHandler
	msg  *communicateData
}

type Room struct {
	activateTime time.Time
	mode         string
	handler      map[uint]*fightingHandler
	isStart      bool
	close        chan bool

	connect chan *fightingHandler
	upload  chan *uploadMessage
	leave   chan *fightingHandler
}

// 房间连接
func (r *Room) attach(handler *fightingHandler) {
	online, ok := r.handler[uint(handler.userId)]
	if !ok {
		bytes, _ := json.Marshal(map[string]any{
			"type": "err",
			"msg":  "not room member",
		})
		handler.msg <- bytes
		return
	}
	if online != nil {
		bytes, _ := json.Marshal(map[string]any{
			"type": "err",
			"msg":  "already online",
		})
		handler.msg <- bytes
		return
	}
	handler.room = r
	r.handler[uint(handler.userId)] = handler
}

// 房间退出
func (r *Room) remove(handler *fightingHandler) {
	oldhandler, ok := r.handler[uint(handler.userId)]
	if !ok {
		return
	}
	if oldhandler != handler {
		return
	}
	r.handler[uint(handler.userId)] = nil
	// 广而告之
	bytes, _ := json.Marshal(map[string]any{
		"type": "quit",
		"msg":  "",
	})
	for _, handler := range r.handler {
		if handler != nil {
			handler.msg <- bytes
		}
	}
}

// 消息传递
func (r *Room) distribute(msg *uploadMessage) {
	bytes, _ := json.Marshal(map[string]any{
		"type": "comm",
		"msg":  msg.msg,
	})
	for _, handler := range r.handler {
		if handler != nil && handler != msg.from {
			handler.msg <- bytes
		}
	}
}

// 游戏启动
func (r *Room) startGame() {
	if !r.isStart {
		allLogin := true
		for _, login := range r.handler {
			if login == nil {
				allLogin = false
				break
			}
		}
		if allLogin {
			// 广播
			bytes, _ := json.Marshal(map[string]any{
				"type": "start",
				"msg":  "",
			})
			for _, handler := range r.handler {
				handler.msg <- bytes
			}
			r.isStart = true
		}
	}
}

// 房间的主循环
func (r *Room) Run() {
	ticker := time.NewTicker(viper.GetDuration("socket.startGameInterval"))
	defer ticker.Stop()
	for {
		select {
		case handler := <-r.connect:
			r.attach(handler)
		case handler := <-r.leave:
			r.remove(handler)
		case msg := <-r.upload:
			r.activateTime = time.Now()
			r.distribute(msg)
		case <-ticker.C:
			r.startGame()
		case <-r.close:
			return
		}
	}
}

type communicateData struct {
	UserId int `json:"user"`
	Score  int `json:"score"`
	Life   int `json:"life"`
}

type fightingHandler struct {
	ws     *websocket.Conn
	userId int
	roomId int
	room   *Room
	msg    chan []byte
}

func (h *fightingHandler) Read() {
	defer h.Close()
	h.ws.SetReadLimit(viper.GetInt64("socket.maxMessageSize"))
	h.ws.SetReadDeadline(time.Now().Add(viper.GetDuration("socket.pongWait")))
	h.ws.SetPongHandler(func(string) error {
		h.ws.SetReadDeadline(time.Now().Add(viper.GetDuration("socket.pongWait")))
		return nil
	})
	for {
		messageType, message, err := h.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Fighting mode websocket read error: %v\n", err)
			}
			break
		}
		if messageType != websocket.TextMessage {
			log.Printf("Fighting mode websocket messageType error: %d\n", messageType)
			break
		}
		data := &communicateData{}
		if err := json.Unmarshal(message, data); err != nil {
			log.Printf("Json unmarshalling error: %v\n", err)
			break
		}
		data.UserId = h.userId
		if h.room != nil {
			h.room.upload <- &uploadMessage{
				from: h,
				msg:  data,
			}
		}
	}
}

func (h *fightingHandler) Write() {
	defer h.Close()
	ticker := time.NewTicker(viper.GetDuration("socket.pingPeriod"))
	defer ticker.Stop()
	for {
		select {
		case message, ok := <-h.msg:
			h.ws.SetWriteDeadline(time.Now().Add(viper.GetDuration("socket.writeWait")))
			if !ok {
				// daemon 关闭了通道
				h.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := h.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Fighting mode websocket write error: %v\n", err)
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				log.Printf("Fighting mode websocket close error: %v\n", err)
				return
			}
		case <-ticker.C:
			h.ws.SetWriteDeadline(time.Now().Add(viper.GetDuration("socket.writeWait")))
			if err := h.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Fighting mode websocket write error: %v\n", err)
				return
			}
		}
	}
}

func (h *fightingHandler) Close() {
	h.ws.Close()
	if h.room != nil {
		h.room.leave <- h
	}
}
