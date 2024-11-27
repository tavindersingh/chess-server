package timer

// Define an interface for notifying players
type PlayerNotifier interface {
	NotifyTimerUpdate(gameID, playerID string, remainingTime float64)
}
