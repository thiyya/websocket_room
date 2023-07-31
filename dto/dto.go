package dto

type RegisterRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

type RegisterResponse struct {
	Id string `json:"id"`
}

type StatsResponse struct {
	RegisteredPlayers int          `json:"registeredPlayers"`
	ActiveRooms       []ActiveRoom `json:"activeRooms"`
}

type ActiveRoom struct {
	Id     string `json:"id"`
	Secret int    `json:"secret"`
}

type Error struct {
	Item    string `json:"items"`
	Message string `json:"message"`
}

type WebSocketRequest struct {
	Cmd string `json:"cmd"`
}

type JoinRequest struct {
	Cmd string `json:"cmd"`
	Id  string `json:"id"`
}

type GuessRequest struct {
	Cmd    string `json:"cmd"`
	Id     string `json:"id"`
	RoomId string `json:"roomId"`
	Data   int    `json:"data"`
}

type WebsocketCommandResponse struct {
	Cmd   string `json:"cmd,omitempty"`
	Reply string `json:"reply,omitempty"`
	Error string `json:"error,omitempty"`
}

type WebsocketEventResponse struct {
	Event    string    `json:"event,omitempty"`
	Room     string    `json:"room,omitempty"`
	Secret   int       `json:"secret,omitempty"`
	Rankings []Ranking `json:"rankings,omitempty"`
}

type Ranking struct {
	Player      string `json:"player"`
	Rank        int    `json:"rank"`
	Guess       int    `json:"guess"`
	DeltaTrophy int    `json:"deltaTrophy"`
}
