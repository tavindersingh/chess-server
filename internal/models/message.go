package models

type Message struct {
	Event string `json:"event"`
	// Payload *Payload `json:"payload,omitempty"`
}

type TimerUpdatePayload struct {
	GameId          string
	Player1Id       string
	Player2Id       string
	Player1TimeLeft int64
	Player2TimeLeft int64
}

type TimeoutPayload struct {
	GameId   string
	WinnerId string
}
