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
	checkDeadInterval = time.Second * 2
)

type GamingController struct {
	handler0  *gamingHandler
	upload0   chan []byte
	download0 chan []byte
	exit0     bool
	handler1  *gamingHandler
	upload1   chan []byte
	download1 chan []byte
	exit1     bool
}

type gamingHandler struct {
	user     *models.User
	ws       *websocket.Conn
	ctrl     *GamingController
	upload   chan []byte
	download chan []byte
	exit     bool
}

func (p *pairingHandler) migrate(ctrl *GamingController) *gamingHandler {
	p.migrated = true
	p.Close()
	return &gamingHandler{
		user:     p.user,
		ws:       p.ws,
		ctrl:     ctrl,
		upload:   make(chan []byte),
		download: make(chan []byte),
		exit:     false,
	}
}

func Fighting(handler0 *pairingHandler, handler1 *pairingHandler) {
	ctrl := &GamingController{}
	newHandler0 := handler0.migrate(ctrl)
	ctrl.handler0 = newHandler0
	ctrl.upload0 = newHandler0.upload
	ctrl.download0 = newHandler0.download
	ctrl.exit0 = newHandler0.exit

	newHandler1 := handler1.migrate(ctrl)
	ctrl.handler1 = newHandler1
	ctrl.upload1 = newHandler1.upload
	ctrl.download1 = newHandler1.download
	ctrl.exit1 = newHandler1.exit

	startMsg, err := json.Marshal(map[string]any{
		"type": "game_start",
	})
	if err != nil {
		log.Printf("Fighting start message marshal error: %v\n", err)
		ctrl.Shutdown()
		return
	}
	go newHandler0.Read()
	go newHandler0.Write()
	go newHandler1.Read()
	go newHandler1.Write()
	ctrl.download0 <- startMsg
	ctrl.download1 <- startMsg
	ctrl.Run()
}

func (con *GamingController) Shutdown() {
	con.handler0.Close()
	con.handler1.Close()
	close(con.upload0)
	close(con.download0)
	close(con.upload1)
	close(con.download1)
}

func (con *GamingController) Run() {
	ticker := time.NewTicker(checkDeadInterval)
	defer ticker.Stop()
	for {
		select {
		case data := <-con.upload0:
			if !con.exit1 {
				con.download1 <- data
			}
		case data := <-con.upload1:
			if !con.exit0 {
				con.download0 <- data
			}
		case <-ticker.C:
			if con.exit0 && con.exit1 {
				con.Shutdown()
				return
			}
		}
	}
}

func (h *gamingHandler) Read() {
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
		h.ctrl.upload0 <- []byte{} // todo
	}
	h.Close()
}

func (h *gamingHandler) Write() {
	ticker := time.NewTicker(viper.GetDuration("socket.pingPeriod"))
	defer func() {
		ticker.Stop()
		h.Close()
	}()
	for {
		select {
		case message, ok := <-h.download:
			h.ws.SetWriteDeadline(time.Now().Add(viper.GetDuration("socket.writeWait")))
			if !ok {
				// ctrl 关闭了通道
				h.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := h.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			h.ws.SetWriteDeadline(time.Now().Add(viper.GetDuration("socket.writeWait")))
			if err := h.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

func (h *gamingHandler) Close() {
	h.ws.Close()
	h.exit = true
}
