package repo

import (
	"github.com/google/uuid"
	"net/http"
	"rooms/global"
	"rooms/model"
	"sync"
)

type WriteRepo interface {
	Update(p *model.Player) error
	Register(name string) (string, error)
	Join(p *model.Player) error
	CreateRoom(r *model.Room) error
	RemoveFromWaitingList(id string) error
}

type ReadRepo interface {
	GetAllPlayers() map[string]*model.Player
	GetPlayerById(id string) (*model.Player, error)
	GetWaitingList() map[string]*model.Player
	GetPlayerByNickName(nickName string) (*model.Player, error)
	GetAllRooms() map[string]*model.Room
	GetRoomById(id string) (*model.Room, error)
}

type Repo interface {
	WriteRepo
	ReadRepo
}

var _ Repo = (*repo)(nil)

type repo struct {
	rooms       map[string]*model.Room
	players     map[string]*model.Player
	waitingList map[string]*model.Player
	mutex       sync.RWMutex
}

func NewRepository(options ...func(*repo)) *repo {
	ar := &repo{}
	for _, o := range options {
		o(ar)
	}
	return ar
}

func WithPlayers(players map[string]*model.Player) func(*repo) {
	return func(s *repo) {
		s.players = players
	}
}

func WithRooms(rooms map[string]*model.Room) func(*repo) {
	return func(s *repo) {
		s.rooms = rooms
	}
}

func WithWaitingList(waitingList map[string]*model.Player) func(*repo) {
	return func(s *repo) {
		s.waitingList = waitingList
	}
}
func (a *repo) Register(nickName string) (string, error) {
	uuid := uuid.New().String()
	a.mutex.Lock()
	a.players[uuid] = &model.Player{
		ID:       uuid,
		NickName: nickName,
		Score:    0,
		Guess:    -1,
	}
	a.mutex.Unlock()
	return uuid, nil
}

func (a *repo) Update(p *model.Player) error {
	a.mutex.Lock()
	a.players[p.ID] = p
	a.mutex.Unlock()
	return nil
}

func (a *repo) GetAllPlayers() map[string]*model.Player {
	a.mutex.RLock()
	p := a.players
	a.mutex.RUnlock()
	return p
}

func (a *repo) GetPlayerById(id string) (*model.Player, error) {
	a.mutex.RLock()
	p, ok := a.players[id]
	a.mutex.RUnlock()
	if ok {
		return p, nil
	}

	return nil, global.NewError(http.StatusNotFound, global.NotFoundErr, "player not found")
}

func (a *repo) GetRoomById(id string) (*model.Room, error) {
	a.mutex.RLock()
	p, ok := a.rooms[id]
	a.mutex.RUnlock()
	if ok {
		return p, nil
	}

	return nil, global.NewError(http.StatusNotFound, global.NotFoundErr, "room not found")
}

func (a *repo) GetPlayerByNickName(nickName string) (*model.Player, error) {
	p := &model.Player{}
	a.mutex.RLock()
	for _, player := range a.players {
		if player.NickName == nickName {
			p = player
			break
		}
	}
	a.mutex.RUnlock()
	if len(p.ID) > 0 {
		return p, nil
	}

	return nil, global.NewError(http.StatusNotFound, global.NotFoundErr, "player not found")
}

func (a *repo) GetAllRooms() map[string]*model.Room {
	a.mutex.RLock()
	p := a.rooms
	a.mutex.RUnlock()
	return p
}

func (a *repo) Join(p *model.Player) error {
	a.mutex.Lock()
	a.waitingList[p.ID] = p
	a.mutex.Unlock()
	return nil
}
func (a *repo) CreateRoom(r *model.Room) error {
	a.mutex.Lock()
	a.rooms[r.ID] = r
	a.mutex.Unlock()
	return nil
}

func (a *repo) RemoveFromWaitingList(id string) error {
	a.mutex.Lock()
	delete(a.waitingList, id)
	a.mutex.Unlock()
	return nil
}

func (a *repo) GetWaitingList() map[string]*model.Player {
	a.mutex.RLock()
	p := a.waitingList
	a.mutex.RUnlock()
	return p
}
