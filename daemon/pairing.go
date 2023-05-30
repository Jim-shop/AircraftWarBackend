package daemon

import (
	"encoding/json"
	"fmt"
	"imshit/aircraftwar/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type PairingDaemon struct {
	handlers map[*pairingHandler]map[uint]bool
	bind     chan *pairingHandler
	unbind   chan *pairingHandler
	pairing  chan *pairRequest
	room     chan *acquireRoomResponce
}

var pairingDaemon *PairingDaemon = nil

func GetPairingDaemon() *PairingDaemon {
	if pairingDaemon == nil {
		pairingDaemon = &PairingDaemon{
			handlers: make(map[*pairingHandler]map[uint]bool),
			bind:     make(chan *pairingHandler, 16),
			unbind:   make(chan *pairingHandler, 16),
			pairing:  make(chan *pairRequest, 16),
			room:     make(chan *acquireRoomResponce, 16),
		}
	}
	return pairingDaemon
}

// 监视用户上线
func (d *PairingDaemon) add(handler *pairingHandler) {
	for current := range d.handlers {
		if current.user.ID == handler.user.ID {
			d.remove(current)
		}
	}
	d.handlers[handler] = make(map[uint]bool)
}

// 删除用户
func (d *PairingDaemon) remove(handler *pairingHandler) {
	if _, ok := d.handlers[handler]; ok {
		close(handler.msg)
		delete(d.handlers, handler)
	}
	userId := handler.user.ID
	for h, req := range d.handlers {
		if _, ok := req[userId]; ok {
			delete(d.handlers[h], userId)
		}
	}
}

// 广播在线信息
func (d *PairingDaemon) boardcast() {
	for handler := range d.handlers {
		player := []any{}
		for other, req := range d.handlers {
			if handler.user.ID != other.user.ID && handler.mode == other.mode {
				_, requesting := req[handler.user.ID]
				player = append(player, map[string]any{
					"ID":         other.user.ID,
					"name":       other.user.Name,
					"requesting": requesting,
				})
			}
		}
		bytes, err := json.Marshal(map[string]any{
			"type":   "player_list",
			"player": player,
		})
		if err != nil {
			log.Printf("Json Marshalling error: %v\n", err)
			continue
		}
		handler.msg <- bytes
	}
}

// 处理匹配
func (d *PairingDaemon) pair(pairing *pairRequest) bool {
	// 配对请求
	sender := pairing.sender
	// 校验发送者
	selectBy, ok := d.handlers[sender]
	if !ok {
		log.Printf("Pairing sender not found.\n")
		return false
	}
	// 校验是否自配对
	if pairing.TargetID == sender.user.ID {
		return false
	}
	// 检查有无接收者
	var target *pairingHandler = nil
	for handler := range d.handlers {
		if pairing.TargetID == handler.user.ID {
			target = handler
		}
	}
	if target == nil {
		log.Printf("Pairing target not found.\n")
		return false
	}
	// 检查难度是否相同
	if target.mode != sender.mode {
		return false
	}
	// 双方进行配对
	d.handlers[target][pairing.sender.user.ID] = true
	// 判断配对情况
	if _, ok := selectBy[pairing.TargetID]; ok {
		fightingDaemon.request <- &acquireRoomRequest{
			master: d,
			h0:     sender,
			h1:     target,
		}
	}
	return true
}

// 处理房间分配完毕
func (d *PairingDaemon) getRoom(responce *acquireRoomResponce) {
	if responce.err != nil {
		return
	}
	h0 := responce.request.h0
	h1 := responce.request.h1
	defer d.remove(h0)
	defer d.remove(h1)
	bytes, err := json.Marshal(map[string]any{
		"type": "room",
		"room": responce.roomID,
	})
	if err != nil {
		log.Printf("Json Marshalling error: %v\n", err)
		return
	}
	h0.msg <- bytes
	h1.msg <- bytes
}

// 主循环
func (d *PairingDaemon) Run() {
	ticker := time.NewTicker(viper.GetDuration("socket.onlinePushInterval"))
	defer ticker.Stop()
	for {
		select {
		case handler := <-d.bind:
			d.add(handler)
		case handler := <-d.unbind:
			d.remove(handler)
		case pairing := <-d.pairing:
			d.pair(pairing)
		case responce := <-d.room:
			d.getRoom(responce)
		case <-ticker.C:
			d.boardcast()
		}
	}
}

// 公共接口
func (d *PairingDaemon) Bind(ws *websocket.Conn, user *models.User, mode string) {
	handler := &pairingHandler{
		user:     user,
		mode:     mode,
		ws:       ws,
		ctrl:     pairingDaemon,
		msg:      make(chan []byte),
		migrated: false,
	}
	pairingDaemon.bind <- handler
	go handler.Read()
	go handler.Write()
}

type pairRequest struct {
	sender   *pairingHandler
	TargetID uint
}

type pairingHandler struct {
	user     *models.User
	mode     string
	ws       *websocket.Conn
	ctrl     *PairingDaemon
	msg      chan []byte
	migrated bool
}

func (h *pairingHandler) Read() {
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
				log.Printf("Pairing mode websocket read error: %v\n", err)
			}
			break
		}
		if messageType != websocket.TextMessage {
			log.Printf("Pairing mode websocket messageType error: %d\n", messageType)
			break
		}
		var targetID uint
		if n, err := fmt.Sscan(string(message), &targetID); n != 1 || err != nil {
			log.Printf("Pairing mode websocket message error: %v\n", err)
			break
		}
		h.ctrl.pairing <- &pairRequest{
			sender:   h,
			TargetID: targetID,
		}
	}
}

func (h *pairingHandler) Write() {
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
				log.Printf("Pairing mode websocket write error: %v\n", err)
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				log.Printf("Pairing mode websocket close error: %v\n", err)
				return
			}
		case <-ticker.C:
			h.ws.SetWriteDeadline(time.Now().Add(viper.GetDuration("socket.writeWait")))
			if err := h.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Pairing mode websocket write error: %v\n", err)
				return
			}
		}
	}
}

func (h *pairingHandler) Close() {
	h.ws.Close()
	h.ctrl.unbind <- h
}
