package models

type Move struct {
	Move      string `json:"move"`
	PlayerId  string `json:"playerId"`
	TimeTaken int64  `json:"timeTaken"`
}
