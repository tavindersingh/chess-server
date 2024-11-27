package timer

import (
	"sync"
	"time"
)

type TimeUpdateCallback func(
	gameId, playerId string,
	remainingTime time.Duration,
)

type Timer struct {
	GameId        string
	PlayerId      string
	startTime     time.Time
	RemainingTime time.Duration
	Ticker        *time.Ticker
	Paused        bool
	StopChan      chan struct{}
	mu            sync.Mutex
	OnTimeout     func()
	OnTimeUpdate  TimeUpdateCallback
}

func NewTimer(
	gameId, playerId string,
	duration time.Duration,
	onTimeOut func(),
	onTimeUpdate TimeUpdateCallback,
) *Timer {
	return &Timer{
		GameId:        gameId,
		PlayerId:      playerId,
		startTime:     time.Now(),
		RemainingTime: duration,
		Paused:        true,
		Ticker:        nil,
		StopChan:      make(chan struct{}),
		OnTimeout:     onTimeOut,
		OnTimeUpdate:  onTimeUpdate,
	}
}

func (t *Timer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Paused {
		return
	}

	t.Paused = false
	t.Ticker = time.NewTicker(1 * time.Second)

	go t.run()
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Ticker != nil {
		t.Ticker.Stop()
	}
	t.RemainingTime = 0
	close(t.StopChan)
	t.StopChan = make(chan struct{})
}

func (t *Timer) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Paused {
		return
	}

	t.Paused = true
	t.Ticker.Stop()
}

func (t *Timer) run() {
	for {
		select {
		case <-t.Ticker.C:
			t.RemainingTime -= time.Second
			if t.OnTimeUpdate != nil {
				t.OnTimeUpdate(
					t.GameId,
					t.PlayerId,
					t.RemainingTime,
				)
			}

			if t.RemainingTime <= 0 {
				t.Ticker.Stop()
				if t.OnTimeout != nil {
					t.OnTimeout()
				}
				
				t.mu.Unlock()
				return
			}
		case <-t.StopChan:
			t.mu.Lock()
			t.Ticker.Stop()
			t.mu.Unlock()
			return
		}
	}
}

// func (t *Timer) EmitTimeUpdate(
// 	gameId, playerId string,
// 	remaining time.Duration,
// ) {
// 	// Example function to broadcast timer update to the players

// }
