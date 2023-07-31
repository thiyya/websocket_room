package model

type Player struct {
	ID       string
	NickName string
	Score    int
	Guess    int
	Diff     int
	Rank     int
}

type Room struct {
	ID      string
	Players []*Player
	Secret  int
}

type Stats struct {
	RegisteredPlayers int
	ActiveRooms       map[string]*Room
}

type GameResult struct {
	Secret   int
	Rankings []Ranking
}

type Ranking struct {
	Player      Player
	Rank        int
	DeltaTrophy int
}
