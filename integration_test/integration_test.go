package integration_test

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"rooms/global"
	"time"

	"rooms/dto"
	"testing"
)

func TestJoinCommand(t *testing.T) {
	Convey("Join", t, func(c C) {
		ws, s := prepareJoin(c)
		defer ws.Close()
		defer s.Close()
		Convey("Join Successfully", func(c C) {
			joinReq := dto.JoinRequest{
				Cmd: "join",
				Id:  "1",
			}
			b, _ := json.Marshal(joinReq)
			err := ws.WriteMessage(websocket.TextMessage, b)
			c.So(err, ShouldBeNil)

			res := dto.WebsocketCommandResponse{
				Cmd:   "join",
				Reply: "waiting",
			}
			b, _ = json.Marshal(res)
			_, p, err := ws.ReadMessage()
			c.So(err, ShouldBeNil)
			c.So(string(p), ShouldEqual, string(b))
		})

		Convey("Join notRegistered", func(c C) {
			joinReq := dto.JoinRequest{
				Cmd: "join",
				Id:  "4",
			}
			b, _ := json.Marshal(joinReq)
			err := ws.WriteMessage(websocket.TextMessage, b)
			c.So(err, ShouldBeNil)
			res := dto.WebsocketCommandResponse{
				Cmd:   "join",
				Error: global.NotRegistered,
			}
			b, _ = json.Marshal(res)
			_, p, err := ws.ReadMessage()
			c.So(err, ShouldBeNil)
			c.So(string(p), ShouldEqual, string(b))
		})

	})
}
func TestJoinedRoomEvent(t *testing.T) {
	Convey("Join", t, func(c C) {
		ws, s := prepareJoinedRoom(c)
		defer ws.Close()
		defer s.Close()
		Convey("JoinedRoomEvent Successfully", func(c C) {
			joinReq := dto.JoinRequest{
				Cmd: "join",
				Id:  "3",
			}
			b, _ := json.Marshal(joinReq)
			err := ws.WriteMessage(websocket.TextMessage, b)
			c.So(err, ShouldBeNil)
			_, p, err := ws.ReadMessage()
			c.So(err, ShouldBeNil)

			time.Sleep(global.CreateRoomTime * time.Second)
			_, p, err = ws.ReadMessage()
			c.So(err, ShouldBeNil)
			c.So(string(p), ShouldContainSubstring, "room")
		})

	})
}
func TestGuessCommand(t *testing.T) {
	Convey("Guess", t, func(c C) {
		ws, s := prepareGuess(c)
		defer ws.Close()
		defer s.Close()
		Convey("Guess Successfully", func(c C) {
			guessReq := dto.GuessRequest{
				Cmd:    "guess",
				Id:     "2",
				RoomId: "room1",
				Data:   5,
			}
			b, _ := json.Marshal(guessReq)
			err := ws.WriteMessage(websocket.TextMessage, b)
			c.So(err, ShouldBeNil)
			res := dto.WebsocketCommandResponse{
				Cmd:   "guess",
				Reply: "guessReceived",
			}
			b, _ = json.Marshal(res)
			_, p, err := ws.ReadMessage()
			c.So(err, ShouldBeNil)
			c.So(string(p), ShouldEqual, string(b))
		})

		Convey("Guess NotRegistered", func(c C) {
			guessReq := dto.GuessRequest{
				Cmd:    "guess",
				Id:     "5",
				RoomId: "room1",
				Data:   5,
			}
			b, _ := json.Marshal(guessReq)
			err := ws.WriteMessage(websocket.TextMessage, b)
			c.So(err, ShouldBeNil)
			res := dto.WebsocketCommandResponse{
				Cmd:   "guess",
				Error: global.NotRegistered,
			}
			b, _ = json.Marshal(res)
			_, p, err := ws.ReadMessage()
			c.So(err, ShouldBeNil)
			c.So(string(p), ShouldEqual, string(b))
		})

		Convey("Guess NotInRoom", func(c C) {
			guessReq := dto.GuessRequest{
				Cmd:    "guess",
				Id:     "4",
				RoomId: "room1",
				Data:   5,
			}
			b, _ := json.Marshal(guessReq)
			err := ws.WriteMessage(websocket.TextMessage, b)
			c.So(err, ShouldBeNil)
			res := dto.WebsocketCommandResponse{
				Cmd:   "guess",
				Error: global.NotInRoom,
			}
			b, _ = json.Marshal(res)
			_, p, err := ws.ReadMessage()
			c.So(err, ShouldBeNil)
			c.So(string(p), ShouldEqual, string(b))
		})
	})
}
func TestGameOverEvent(t *testing.T) {
	Convey("GameOver", t, func(c C) {
		ws, s := prepareGameOver(c)
		defer ws.Close()
		defer s.Close()
		Convey("GameOver Successfully", func(c C) {
			guessReq := dto.GuessRequest{
				Cmd:    "guess",
				Id:     "3",
				RoomId: "room1",
				Data:   3,
			}
			b, _ := json.Marshal(guessReq)
			err := ws.WriteMessage(websocket.TextMessage, b)
			c.So(err, ShouldBeNil)
			r := dto.WebsocketCommandResponse{
				Cmd:   "guess",
				Reply: "guessReceived",
			}
			b, _ = json.Marshal(r)
			_, p, err := ws.ReadMessage()
			c.So(err, ShouldBeNil)
			c.So(string(p), ShouldEqual, string(b))

			_, p, err = ws.ReadMessage()
			c.So(err, ShouldBeNil)
			res := &dto.WebsocketEventResponse{}
			err = json.Unmarshal(p, res)
			c.So(err, ShouldBeNil)
			expected := &dto.WebsocketEventResponse{
				Event:  "gameOver",
				Secret: 3,
				Rankings: []dto.Ranking{
					{Player: "3", Rank: 1, Guess: 3, DeltaTrophy: 30},
					{Player: "2", Rank: 2, Guess: 4, DeltaTrophy: 20},
					{Player: "1", Rank: 3, Guess: 5, DeltaTrophy: 0},
				},
			}
			ex, _ := json.Marshal(expected)
			c.So(string(p), ShouldEqual, string(ex))
		})

	})
}
