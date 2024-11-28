package game

import "sync"

type CacheGameRepository struct {
	games map[string]*Game
	mu    sync.RWMutex
}

func NewCacheGameRepository() *CacheGameRepository {
	return &CacheGameRepository{
		games: make(map[string]*Game),
	}
}

func (cgr *CacheGameRepository) GetGame(gameId string) (*Game, bool) {
	cgr.mu.RLock()
	defer cgr.mu.RUnlock()
	game, exists := cgr.games[gameId]
	return game, exists
}

func (cgr *CacheGameRepository) AddGame(game *Game) error {
	cgr.mu.Lock()
	defer cgr.mu.Unlock()
	cgr.games[game.GameId] = game
	return nil
}

func (cgr *CacheGameRepository) UpdateGame(gameId string, game *Game) (*Game, bool) {
	cgr.mu.Lock()
	defer cgr.mu.Unlock()

	if _, exists := cgr.games[gameId]; !exists {
		return nil, false
	}

	cgr.games[gameId] = game
	return game, true
}
