package models

type Message struct {
	Event   string   `json:"event"`
	Payload *Payload `json:"payload,omitempty"`
}

type Payload struct {
	Move     string `json:"move"`
	PlayerId string `json:"playerId"`
	GameId   string `json:"gameId"`
}

type MoveMessage struct {
	Event   string       `json:"event"`
	Payload *MovePayload `json:"payload"`
}

type MovePayload struct {
	Event           string `json:"event"`
	GameId          string `json:"gameId"`
	Player1Id       string `json:"player1Id"`
	Player2Id       string `json:"player2Id"`
	LastMove        *Move  `json:"lastMove"`
	CurrentPlayerId string `json:"currentPlayerId"`
}

type TimerUpdateMessage struct {
	Event   string              `json:"event"`
	Payload *TimerUpdatePayload `json:"payload"`
}

type TimerUpdatePayload struct {
	GameId          string `json:"gameId"`
	Player1Id       string `json:"player1Id"`
	Player2Id       string `json:"player2Id"`
	Player1TimeLeft int64  `json:"player1TimeLeft"`
	Player2TimeLeft int64  `json:"player2TimeLeft"`
}

type TimeoutPayload struct {
	GameId   string
	WinnerId string
}
