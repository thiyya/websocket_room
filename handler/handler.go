package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"rooms/dto"
	"rooms/global"
	"rooms/service"
	"sync"
	"time"
)

type Handler struct {
	service *service.Service
	conn    *websocket.Conn
	sync.Mutex
}

func NewHandler(options ...func(*Handler)) *Handler {
	as := &Handler{}
	for _, o := range options {
		o(as)
	}
	return as
}

func WithService(s *service.Service) func(*Handler) {
	return func(h *Handler) {
		h.service = s
	}
}

func (a *Handler) Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := dto.RegisterRequest{}
		json.NewDecoder(r.Body).Decode(&input)
		err := validator.New().Struct(input)
		if err != nil {
			log.Printf("err occurred while parsing register input: %s \n", err.Error())
			err = global.NewError(http.StatusBadRequest, global.InvalidParams, err.Error())
			err.(*global.Error).WriteError(w)
			return
		}
		id, err := a.service.Register(input.Nickname)
		if err != nil {
			log.Printf("err occurred while registering : %s \n", err.Error())
			err.(*global.Error).WriteError(w)

			return
		}
		writeResponse(w, dto.RegisterResponse{Id: id}, http.StatusCreated)
	})
}

func (a *Handler) Stats() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stats, err := a.service.Stats()
		if err != nil {
			log.Printf("err occurred while getting stats : %s \n", err.Error())
			err.(*global.Error).WriteError(w)

			return
		}
		activeRooms := make([]dto.ActiveRoom, 0)
		for _, s := range stats.ActiveRooms {
			activeRooms = append(activeRooms, dto.ActiveRoom{
				Id:     s.ID,
				Secret: s.Secret,
			})
		}
		writeResponse(w, dto.StatsResponse{
			RegisteredPlayers: stats.RegisteredPlayers,
			ActiveRooms:       activeRooms,
		}, http.StatusOK)
	})
}

func (a *Handler) Websocket() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Websocket Connection error : ", err)
			return
		}
		log.Println("Websocket Connection established.")
		a.conn = conn
		defer a.conn.Close()

		clientId := make(chan string)
		gameOver := make(chan bool)

		a.handleJoinedRoomEvent(a.conn, clientId)
		a.handleGameOverEvent(a.conn, gameOver)
		a.handleCommand(a.conn, clientId, gameOver)
	})
}

func (a *Handler) handleCommand(conn *websocket.Conn, clientId chan string, gameOver chan bool) {
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				return
			}
			log.Println("Could not read message from websocket, error : ", err)
			res := dto.WebsocketCommandResponse{
				Error: global.InvalidRequest,
			}
			b, err := json.Marshal(res)
			a.Lock()
			if err = conn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Println("Could not write message to websocket, error", err)
			}
			a.Unlock()
			continue
		}
		commandRequest := &dto.WebSocketRequest{}
		err = json.Unmarshal(b, commandRequest)
		if err != nil {
			log.Println("Could not unmarshal the read message of websocket, error", err)
			res := dto.WebsocketCommandResponse{
				Error: global.InvalidRequest,
			}
			b, err := json.Marshal(res)
			a.Lock()
			if err = conn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Println("Could not write message to websocket, error", err)
			}
			a.Unlock()
			continue
		}

		switch commandRequest.Cmd {
		case "join":
			joinRequest := &dto.JoinRequest{}
			err = json.Unmarshal(b, joinRequest)
			if err != nil {
				log.Println("Could not unmarshal the join message of websocket, error", err)
				res := dto.WebsocketCommandResponse{
					Cmd:   commandRequest.Cmd,
					Error: global.InvalidParams,
				}
				b, err := json.Marshal(res)
				a.Lock()
				if err = conn.WriteMessage(websocket.TextMessage, b); err != nil {
					log.Println("Could not write message to websocket, error", err)
				}
				a.Unlock()
				continue
			}
			err = a.service.Join(joinRequest.Id)
			if err != nil {
				log.Println("Err occurred while join process, error", err)
				res := dto.WebsocketCommandResponse{
					Cmd:   commandRequest.Cmd,
					Error: err.Error(),
				}
				b, err := json.Marshal(res)
				a.Lock()
				if err = conn.WriteMessage(websocket.TextMessage, b); err != nil {
					log.Println("Could not write message to websocket, error", err)
				}
				a.Unlock()
				continue
			}
			res := dto.WebsocketCommandResponse{
				Cmd:   commandRequest.Cmd,
				Reply: "waiting",
			}
			databytes, err := json.Marshal(res)
			a.Lock()
			if err = conn.WriteMessage(websocket.TextMessage, databytes); err != nil {
				log.Println("Could not write message to websocket, error", err)
				a.Unlock()
				continue
			}
			a.Unlock()
			clientId <- joinRequest.Id
		case "guess":
			guessRequest := &dto.GuessRequest{}
			err = json.Unmarshal(b, guessRequest)
			if err != nil {
				log.Println("Could not unmarshal the guess message of websocket, error", err)
				res := dto.WebsocketCommandResponse{
					Cmd:   commandRequest.Cmd,
					Error: global.InvalidParams,
				}
				b, err := json.Marshal(res)
				a.Lock()
				if err = conn.WriteMessage(websocket.TextMessage, b); err != nil {
					log.Println("Could not write message to websocket, error", err)
				}
				a.Unlock()
				continue
			}
			err = a.service.Guess(guessRequest.Id, guessRequest.RoomId, guessRequest.Data)
			if err != nil {
				log.Println("Err occurred while guess process, error", err)
				res := dto.WebsocketCommandResponse{
					Cmd:   commandRequest.Cmd,
					Error: err.Error(),
				}
				b, err := json.Marshal(res)
				a.Lock()
				if err = conn.WriteMessage(websocket.TextMessage, b); err != nil {
					log.Println("Could not write message to websocket, error", err)
				}
				a.Unlock()
				continue
			}
			res := dto.WebsocketCommandResponse{
				Cmd:   commandRequest.Cmd,
				Reply: "guessReceived",
			}
			databytes, err := json.Marshal(res)
			a.Lock()
			if err = conn.WriteMessage(websocket.TextMessage, databytes); err != nil {
				log.Println("Could not write message to websocket, error", err)
				a.Unlock()
				continue
			}
			a.Unlock()
			allGuessDone := a.service.AllGuessDone(guessRequest.RoomId)
			if allGuessDone {
				gameOver <- allGuessDone
			}
		default:
			log.Println("Not supported command, ", commandRequest.Cmd)
		}
	}
}

func (a *Handler) handleJoinedRoomEvent(conn *websocket.Conn, clientId chan string) {
	tickerCreateRoomTime := time.NewTicker(global.CreateRoomTime * time.Second)
	done := make(chan bool)
	go func() {
		id := <-clientId
		for {
			select {
			case <-done:
				return
			case <-tickerCreateRoomTime.C:
				a.Lock()
				rooms := a.service.CreateRooms()
				a.Unlock()
				roomId := ""
				for _, room := range rooms {
					for _, player := range room.Players {
						if player.ID == id {
							roomId = room.ID
							break
						}
					}
				}

				if len(roomId) > 0 {
					res := &dto.WebsocketEventResponse{
						Event: "joinedRoom",
						Room:  roomId,
					}
					databytes, err := json.Marshal(res)
					a.Lock()
					if err = conn.WriteMessage(websocket.TextMessage, databytes); err != nil {
						log.Println("Could not write message to websocket, error", err)
						a.Unlock()
						return
					}
					a.Unlock()
					done <- true
				}
			}
		}
	}()
}

func (a *Handler) handleGameOverEvent(conn *websocket.Conn, gameOver chan bool) {
	go func() {
		<-gameOver
		a.service.GameOver()
		gameResults := a.service.GetGameResults()
		ranking := make([]dto.Ranking, 0)
		for _, r := range gameResults.Rankings {
			ranking = append(ranking, dto.Ranking{
				Player:      r.Player.ID,
				Rank:        r.Rank,
				Guess:       r.Player.Guess,
				DeltaTrophy: r.DeltaTrophy,
			})
		}
		res := &dto.WebsocketEventResponse{
			Event:    "gameOver",
			Secret:   gameResults.Secret,
			Rankings: ranking,
		}
		databytes, err := json.Marshal(res)
		a.Lock()
		if err = conn.WriteMessage(websocket.TextMessage, databytes); err != nil {
			log.Println("Could not write message to websocket, error", err)
			a.Unlock()
			return
		}
		a.Unlock()
	}()
}

func writeResponse(w http.ResponseWriter, v any, responseCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	b, _ := json.Marshal(v)
	w.Write(b)
}
