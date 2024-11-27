package timer

type GameTimers struct {
	Player1Timer *Timer
	Player2Timer *Timer
}

type TimerRepository interface {
	AddGameTimers(gameID string, timers *GameTimers)
	GetGameTimers(gameID string) (*GameTimers, bool)
	RemoveGameTimers(gameID string)
}
