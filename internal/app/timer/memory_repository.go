package timer

import "sync"

type InMemoryTimerRepository struct {
	timers map[string]*GameTimers
	mu     sync.RWMutex
}

func NewInMemoryTimerRepository() *InMemoryTimerRepository {
	return &InMemoryTimerRepository{
		timers: make(map[string]*GameTimers),
	}
}

func (r *InMemoryTimerRepository) AddGameTimers(gameID string, timers *GameTimers) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.timers[gameID] = timers
}

func (repo *InMemoryTimerRepository) GetGameTimers(gameID string) (*GameTimers, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	timers, exists := repo.timers[gameID]
	return timers, exists
}

func (repo *InMemoryTimerRepository) RemoveGameTimers(gameID string) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	delete(repo.timers, gameID)
}
