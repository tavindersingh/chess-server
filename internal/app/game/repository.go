package game

type GameRepository interface {
	AddGame(game *Game) error
	GetGame(gameId string) (*Game, bool)
	UpdateGame(gameId string, game *Game) (*Game, bool)
}
