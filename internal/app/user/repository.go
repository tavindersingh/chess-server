package user

import "github.com/google/uuid"

type UserRepository interface {
	GetUser(id string) (User, error)
	CreateUser() (User, error)
}

type InMemoryUserRepository struct {
	Users map[string]User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		Users: make(map[string]User),
	}
}

func (ur *InMemoryUserRepository) GetUser(id string) (User, error) {
	return ur.Users[id], nil
}

func (ur *InMemoryUserRepository) CreateUser() (User, error) {
	user := &User{
		Id: "anon-" + uuid.NewString(),
	}
	ur.Users[user.Id] = *user

	return *user, nil
}
