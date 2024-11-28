package auth

import (
	"tavinder/chess-server/internal/app/user"
)

type authService struct {
	jwtManager     *JwtManager
	userRepository user.UserRepository
}

type AuthServiceImpl interface {
	AnonymousLogin() (string, error)
}

func NewAuthService(jwtManager *JwtManager, userRepository user.UserRepository) AuthServiceImpl {
	return &authService{
		jwtManager:     jwtManager,
		userRepository: userRepository,
	}
}

func (as *authService) AnonymousLogin() (string, error) {
	user, err := as.userRepository.CreateUser()
	if err != nil {
		return "", err
	}

	return as.jwtManager.GenerateToken(user.Id)
}
