package service

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"rooms/global"
	"rooms/model"
	"rooms/repo"
	"sort"
	"time"
)

type Service struct {
	repo repo.Repo
}

func NewService(options ...func(*Service)) *Service {
	as := &Service{}
	for _, o := range options {
		o(as)
	}
	return as
}

func WithRepo(r repo.Repo) func(*Service) {
	return func(s *Service) {
		s.repo = r
	}
}

func (a *Service) Register(nickName string) (string, error) {
	existedPlayer, err := a.repo.GetPlayerByNickName(nickName)
	if err == nil {
		return existedPlayer.ID, nil
	}

	return a.repo.Register(nickName)
}

func (a *Service) Stats() (model.Stats, error) {
	res := model.Stats{}
	res.RegisteredPlayers = len(a.repo.GetAllPlayers())
	res.ActiveRooms = a.repo.GetAllRooms()
	return res, nil
}

func (a *Service) Join(id string) error {
	p, err := a.repo.GetPlayerById(id)
	if err != nil {
		return fmt.Errorf("%s", global.NotRegistered)
	}
	err = a.repo.Join(p)
	if err != nil {
		return err
	}
	return nil
}

func (a *Service) Guess(id string, roomId string, guess int) error {
	_, err := a.repo.GetPlayerById(id)
	if err != nil {
		return fmt.Errorf("%s", global.NotRegistered)
	}
	r, err := a.repo.GetRoomById(roomId)
	if err != nil {
		return fmt.Errorf("%s", global.NotFoundErr)
	}
	p := &model.Player{}
	for _, player := range r.Players {
		if player.ID == id {
			p = player
			break
		}
	}
	if len(p.ID) == 0 {
		return fmt.Errorf("%s", global.NotInRoom)
	}
	p.Guess = guess
	p.Diff = abs(r.Secret - guess)
	err = a.repo.Update(p)
	if err != nil {
		return err
	}
	return nil
}

func (a *Service) CreateRooms() map[string]*model.Room {
	rand.Seed(time.Now().UnixNano())
	waitingList := a.repo.GetWaitingList()
	p := make([]*model.Player, 0)
	for _, player := range waitingList {
		p = append(p, player)
	}
	for i := 0; i < len(waitingList)-2; i = i + 3 {
		a.repo.CreateRoom(&model.Room{
			ID:      uuid.New().String(),
			Players: p[i : i+3],
			Secret:  rand.Intn(10) + 1,
		})
		a.repo.RemoveFromWaitingList(p[i].ID)
		a.repo.RemoveFromWaitingList(p[i+1].ID)
		a.repo.RemoveFromWaitingList(p[i+2].ID)
	}

	return a.repo.GetAllRooms()
}

func (a *Service) GetGameResults(roomId string) model.GameResult {
	res := model.GameResult{}
	for _, r := range a.repo.GetAllRooms() {
		if roomId == r.ID {
			res.Secret = r.Secret
			rankings := make([]model.Ranking, 0)
			for _, player := range r.Players {
				rankings = append(rankings, model.Ranking{
					Player:      *player,
					Rank:        player.Rank,
					DeltaTrophy: player.Score,
				})
			}
			res.Rankings = rankings
			break
		}
	}
	return res
}
func (a *Service) GameOver(roomId string) {
	rooms := a.repo.GetAllRooms()
	for _, room := range rooms {
		if room.ID == roomId {
			sort.Sort(sortByDiff(room.Players))
			room.Players[0].Score = global.WinnerPrize
			room.Players[0].Rank = 1
			room.Players[1].Score = global.SecondPrize
			room.Players[1].Rank = 2
			room.Players[2].Score = global.Loser
			room.Players[2].Rank = 3
			break
		}
	}
}

func (a *Service) AllGuessDone(roomId string) bool {
	room, _ := a.repo.GetRoomById(roomId)
	for _, player := range room.Players {
		if player.Guess == -1 {
			return false
		}
	}
	return true
}

type sortByDiff []*model.Player

func (e sortByDiff) Len() int {
	return len(e)
}

func (e sortByDiff) Less(i, j int) bool {
	return e[i].Diff < e[j].Diff
}

func (e sortByDiff) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
