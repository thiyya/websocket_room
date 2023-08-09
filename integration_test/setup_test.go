package integration_test

import (
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"rooms/handler"
	"rooms/model"
	"rooms/repo"
	"rooms/service"
	"strings"
)

func prepareJoin(c C) (*websocket.Conn, *httptest.Server) {
	players := map[string]*model.Player{}
	players["1"] = &model.Player{
		ID:       "1",
		NickName: "a",
	}
	players["2"] = &model.Player{
		ID:       "2",
		NickName: "b",
	}
	players["3"] = &model.Player{
		ID:       "3",
		NickName: "c",
	}
	repo := repo.NewRepository(
		repo.WithPlayers(players),
		repo.WithRooms(map[string]*model.Room{}),
		repo.WithWaitingList(map[string]*model.Player{}),
	)
	serv := service.NewService(service.WithRepo(repo))
	handler := handler.NewHandler(handler.WithService(serv))
	s := httptest.NewServer(http.HandlerFunc(handler.Websocket().ServeHTTP))
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	c.So(err, ShouldBeNil)
	return ws, s
}
func prepareJoinedRoom(c C) (*websocket.Conn, *httptest.Server) {
	p1 := &model.Player{
		ID:       "1",
		NickName: "a",
	}
	p2 := &model.Player{
		ID:       "2",
		NickName: "b",
	}
	p3 := &model.Player{
		ID:       "3",
		NickName: "c",
	}
	players := map[string]*model.Player{}
	players["1"] = p1
	players["2"] = p2
	players["3"] = p3
	wl := map[string]*model.Player{}
	wl["1"] = p1
	wl["2"] = p2
	repo := repo.NewRepository(
		repo.WithPlayers(players),
		repo.WithRooms(map[string]*model.Room{}),
		repo.WithWaitingList(wl),
	)
	serv := service.NewService(service.WithRepo(repo))
	handler := handler.NewHandler(handler.WithService(serv))
	s := httptest.NewServer(http.HandlerFunc(handler.Websocket().ServeHTTP))
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	c.So(err, ShouldBeNil)
	return ws, s
}
func prepareGuess(c C) (*websocket.Conn, *httptest.Server) {
	p1 := &model.Player{
		ID:       "1",
		NickName: "a",
	}
	p2 := &model.Player{
		ID:       "2",
		NickName: "b",
	}
	p3 := &model.Player{
		ID:       "3",
		NickName: "c",
	}
	p4 := &model.Player{
		ID:       "4",
		NickName: "d",
	}
	players := map[string]*model.Player{}
	players["1"] = p1
	players["2"] = p2
	players["3"] = p3
	players["4"] = p4
	rooms := map[string]*model.Room{}
	rooms["room1"] = &model.Room{
		ID:      "room1",
		Players: []*model.Player{p1, p2, p3},
		Secret:  0,
	}
	repo := repo.NewRepository(
		repo.WithPlayers(players),
		repo.WithRooms(rooms),
		repo.WithWaitingList(map[string]*model.Player{}),
	)
	serv := service.NewService(service.WithRepo(repo))
	handler := handler.NewHandler(handler.WithService(serv))
	s := httptest.NewServer(http.HandlerFunc(handler.Websocket().ServeHTTP))
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	c.So(err, ShouldBeNil)
	return ws, s
}
func prepareGameOver(c C) (*websocket.Conn, *httptest.Server) {
	p1 := &model.Player{
		ID:       "1",
		NickName: "a",
		Score:    0,
		Guess:    5,
		Diff:     2,
		Rank:     0,
	}
	p2 := &model.Player{
		ID:       "2",
		NickName: "b",
		Score:    0,
		Guess:    4,
		Diff:     1,
		Rank:     0,
	}
	p3 := &model.Player{
		ID:       "3",
		NickName: "c",
	}

	players := map[string]*model.Player{}
	players["1"] = p1
	players["2"] = p2
	players["3"] = p3
	rooms := map[string]*model.Room{}
	rooms["room1"] = &model.Room{
		ID:      "room1",
		Players: []*model.Player{p1, p2, p3},
		Secret:  3,
	}
	repo := repo.NewRepository(
		repo.WithPlayers(players),
		repo.WithRooms(rooms),
		repo.WithWaitingList(map[string]*model.Player{}),
	)
	serv := service.NewService(service.WithRepo(repo))
	handler := handler.NewHandler(handler.WithService(serv))
	s := httptest.NewServer(http.HandlerFunc(handler.Websocket().ServeHTTP))
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	c.So(err, ShouldBeNil)
	return ws, s
}
