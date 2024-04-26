package service

import (
	"sync"

	"github.com/swarit-pandey/e-commerce/user/pkg/repository"
)

type inMemory struct {
	mu    sync.Mutex
	store map[any]*repository.User
}

func newInMemory() *inMemory {
	return &inMemory{
		store: make(map[any]*repository.User),
	}
}

func (im *inMemory) setInMemory(user *repository.User) {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.store[user.Username] = user
}

func (im *inMemory) getFromInMem(key any) *repository.User {
	im.mu.Lock()
	defer im.mu.Unlock()

	user, ok := im.store[key]
	if !ok {
		return nil
	}
	return user
}

func (im *inMemory) getBatchInMem() []*repository.User {
	im.mu.Lock()
	defer im.mu.Unlock()

	users := make([]*repository.User, 0, len(im.store))
	for _, user := range im.store {
		users = append(users, user)
	}

	return users
}
