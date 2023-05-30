package socket

import (
	"encoding/json"
	"fmt"
	"imshit/aircraftwar/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const (
	onlinePushInterval = time.Second * 2
)

type pairRequest struct {
	sender   *pairingHandler
	TargetID uint
}

type PairingController struct {
	handlers map[*pairingHandler]map[uint]bool
	bind     chan *pairingHandler
	unbind   chan *pairingHandler
	pairing  chan pairRequest
}

type pairingHandler struct {
	user     *models.User
	ws       *websocket.Conn
	ctrl     *PairingController
	msg      chan []byte
	migrated bool
}

type onlineInfo struct {
	ID         uint
	Name       string
	Requesting bool
}

// 单例
var pairingController = &PairingController{
	handlers: make(map[*pairingHandler]map[uint]bool),
	bind:     make(chan *pairingHandler),
	unbind:   make(chan *pairingHandler),
	pairing:  make(chan pairRequest),
}

func GetPairingController() *PairingController {
	return pairingController
}

func (con *PairingController) Run() {
	ticker := time.NewTicker(onlinePushInterval)
	defer ticker.Stop()
	for {
		select {
		case handler := <-con.bind:
			// 监视用户上线
			for exist := range con.handlers {
				if exist.user.ID == handler.user.ID {
					con.unbind <- exist
				}
			}
			con.handlers[handler] = make(map[uint]bool)
		case handler := <-con.unbind:
			// 监视用户下线
			if _, ok := con.handlers[handler]; ok {
				close(handler.msg)
				delete(con.handlers, handler)
			}
			userId := handler.user.ID
			for h, req := range con.handlers {
				if _, ok := req[userId]; ok {
					delete(con.handlers[h], userId)
				}
			}
		case pairing := <-con.pairing:
			// 配对请求
			sender := pairing.sender
			// 校验发送者
			selectBy, ok := con.handlers[sender]
			if !ok {
				log.Printf("Pairing sender not found.\n")
				continue
			}
			// 校验接收者
			if pairing.TargetID == sender.user.ID {
				log.Printf("Pairing with oneself is not allow.\n")
				continue
			}
			var target *pairingHandler = nil
			for handler := range con.handlers {
				if pairing.TargetID == handler.user.ID {
					target = handler
				}
			}
			if target == nil {
				log.Printf("Pairing target not found.\n")
				continue
			}
			// 判断配对情况
			if _, ok := selectBy[pairing.TargetID]; ok {
				// 我方同意他方
				go Fighting(sender, target)
			} else {
				// 我方请求他方
				con.handlers[target][pairing.sender.user.ID] = true
			}
		case <-ticker.C:
			// 广播在线消息
			allOnline := []onlineInfo{}
			for handler := range con.handlers {
				allOnline = append(allOnline, onlineInfo{
					ID:   handler.user.ID,
					Name: handler.user.Name,
				})
			}
			for handler, req := range con.handlers {
				otherOnline := []onlineInfo{}
				for _, online := range allOnline {
					_, requesting := req[online.ID]
					if online.ID != handler.user.ID {
						otherOnline = append(otherOnline, onlineInfo{
							ID:         online.ID,
							Name:       online.Name,
							Requesting: requesting,
						})
					}
				}
				bytes, err := json.Marshal(map[string]any{
					"type":   "player_list",
					"player": otherOnline,
				})
				if err != nil {
					log.Printf("Json Marshalling error: %v\n", err)
					handler.Close()
					continue
				}
				handler.msg <- bytes
			}
		}
	}
}

func Pairing(ws *websocket.Conn, user *models.User) {
	handler := &pairingHandler{
		user:     user,
		ws:       ws,
		ctrl:     pairingController,
		msg:      make(chan []byte),
		migrated: false,
	}
	pairingController.bind <- handler
	go handler.Read()
	go handler.Write()
}

func (h *pairingHandler) Read() {
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
		h.ctrl.pairing <- pairRequest{
			sender:   h,
			TargetID: targetID,
		}
	}
	h.Close()
}

func (h *pairingHandler) Write() {
	ticker := time.NewTicker(viper.GetDuration("socket.pingPeriod"))
	defer func() {
		ticker.Stop()
		h.Close()
	}()
	for {
		select {
		case message, ok := <-h.msg:
			h.ws.SetWriteDeadline(time.Now().Add(viper.GetDuration("socket.writeWait")))
			if !ok {
				// ctrl 关闭了通道
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
	h.ctrl.unbind <- h
	if !h.migrated {
		h.ws.Close()
	}
}
