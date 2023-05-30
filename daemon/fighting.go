package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"imshit/aircraftwar/models"

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
	// TODO
}

var fightingDaemon *FightingDaemon = nil

func GetFightingDaemon() *FightingDaemon {
	if fightingDaemon == nil {
		fightingDaemon = &FightingDaemon{
			rooms:   make(map[int]*Room),
			request: make(chan *acquireRoomRequest, 16),
			connect: make(chan *fightingHandler, 16),
			// TODO
		}
	}
	return fightingDaemon
}

// 获取房间
func (d *FightingDaemon) acquire(request *acquireRoomRequest) {
	for i := 0; i < viper.GetInt("socket.maxRoomNum"); i++ {
		_, ok := d.rooms[i]
		if !ok {
			d.rooms[i] = &Room{
				// TODO
				id0: request.h0.user.ID,
				id1: request.h1.user.ID,
			}
			request.master.room <- &acquireRoomResponce{
				request: request,
				roomID:  i,
				err:     nil,
			}
			return
		}
	}
	request.master.room <- &acquireRoomResponce{
		request: request,
		roomID:  0,
		err:     errors.New("exceed max room num"),
	}
}

// 清理无用房间
func (d *FightingDaemon) clean() {

}

// 主循环
func (d *FightingDaemon) Run() {
	ticker := time.NewTicker(viper.GetDuration("socket.cleanRoomInterval"))
	defer ticker.Stop()
	for {
		select {
		case request := <-d.request:
			d.acquire(request)
		case <-ticker.C:
			d.clean()
		}
	}
}

type Room struct {
	id0 uint
	id1 uint
}

type GameDaemon struct {
	handler0  *fightingHandler
	upload0   chan []byte
	download0 chan []byte
	exit0     chan bool
	isExit0   bool
	handler1  *fightingHandler
	upload1   chan []byte
	download1 chan []byte
	exit1     chan bool
	isExit1   bool
	isSilent  bool
}

type fightingHandler struct {
	user     *models.User
	ws       *websocket.Conn
	ctrl     *GameDaemon
	upload   chan []byte
	download chan []byte
	exit     chan bool
}

func (p *pairingHandler) migrate(ctrl *GameDaemon) *fightingHandler {
	p.migrated = true
	handler := &fightingHandler{
		user:     p.user,
		ws:       p.ws,
		ctrl:     ctrl,
		upload:   make(chan []byte),
		download: make(chan []byte),
		exit:     make(chan bool),
	}
	p.Close()
	return handler
}

func Fighting(handler0 *pairingHandler, handler1 *pairingHandler) {
	ctrl := &GameDaemon{
		isExit0: false,
		isExit1: false,
	}

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

	go ctrl.Run()
	go newHandler0.Read()
	go newHandler0.Write()
	go newHandler1.Read()
	go newHandler1.Write()

	go ctrl.SendStart()
}

func (con *GameDaemon) SendStart() {
	startMsg, err := json.Marshal(map[string]any{
		"type": "game_start",
	})
	if err != nil {
		log.Printf("Fighting start message marshal error: %v\n", err)
		con.Shutdown()
		return
	}
	con.download0 <- startMsg
	con.download1 <- startMsg
}

func (con *GameDaemon) Shutdown() {
	con.handler0.Close()
	con.handler1.Close()
	close(con.upload0)
	close(con.download0)
	close(con.upload1)
	close(con.download1)
}

func (con *GameDaemon) Run() {
	ticker := time.NewTicker(viper.GetDuration("socket.checkDeadInterval"))
	defer ticker.Stop()
	for {
		select {
		case data := <-con.upload0:
			con.isSilent = false
			if !con.isExit1 {
				con.download1 <- data
			}
		case data := <-con.upload1:
			con.isSilent = false
			if !con.isExit0 {
				con.download0 <- data
			}
		case <-con.exit0:
			con.isSilent = false
			con.isExit0 = true
		case <-con.exit1:
			con.isSilent = false
			con.isExit1 = true
		case <-ticker.C:
			if (con.isExit0 && con.isExit1) || con.isSilent {
				con.Shutdown()
				return
			}
			con.isSilent = true
		}
	}
}

func (h *fightingHandler) Read() {
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
		var targetID uint
		if n, err := fmt.Sscan(string(message), &targetID); n != 1 || err != nil {
			log.Printf("Fighting mode websocket message error: %v\n", err)
			break
		}
		h.ctrl.upload0 <- []byte{} // todo
	}
	h.Close()
}

func (h *fightingHandler) Write() {
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
				return
			}
		}
	}

}

func (h *fightingHandler) Close() {
	h.ws.Close()
	h.exit <- true
}
