package quiz

import (
	"errors"
	"sync"
)

type usersMap struct {
	container map[string]*User
	mu        sync.RWMutex
}

func newUsersMap() *usersMap {
	return &usersMap{
		container: make(map[string]*User),
		mu:        sync.RWMutex{},
	}
}

func (umap *usersMap) addUser(u *User) {
	umap.mu.Lock()
	umap.container[u.id] = u
	umap.mu.Unlock()
}

func (umap *usersMap) getUser(userID string) *User {
	umap.mu.RLock()
	user, ok := umap.container[userID]
	umap.mu.RUnlock()

	if !ok || user == nil {
		return nil
	}
	return user
}

func (umap *usersMap) isAllReady() bool {
	c := 0
	umap.mu.Lock()
	for _, u := range umap.container {
		if u.state == userStateReadyToPlay {
			c++
		}
	}

	isAllReady := c == len(umap.container)
	umap.mu.Unlock()

	return isAllReady
}
func (umap *usersMap) sendMessageToUser(userID string, msg serverMessage) error {
	umap.mu.RLock()
	defer umap.mu.RUnlock()
	user, ok := umap.container[userID]
	if !ok {
		return errUserNotFound
	}

	return user.sendMessageToUser(msg)
}

func (umap *usersMap) sendMessageForEachUser(msg serverMessage) error {
	umap.mu.RLock()
	defer umap.mu.RUnlock()

	var anyErr error

	for _, u := range umap.container {
		err := u.sendMessageToUser(msg)
		if err != nil {
			anyErr = err
		}
	}

	if anyErr != nil {
		return anyErr
	}
	return nil
}

// on gameover
func (umap *usersMap) SendResultsForEachUser() error {
	umap.mu.RLock()
	defer umap.mu.RUnlock()

	var anyErr error

	// collect all users score
	usersScore := make(map[string]int)
	for _, u := range umap.container {
		usersScore[u.id] = u.score
	}

	msg := serverMessage{
		MessageType: serverMessageTypeGameOver,
		Message:     usersScore,
	}
	for _, u := range umap.container {
		err := u.sendMessageToUser(msg)
		if err != nil {
			anyErr = err
		}
	}

	if anyErr != nil {
		return anyErr
	}
	return nil
}

func (umap *usersMap) setUsersToBuffered() {
	// reset all users
	umap.mu.Lock()
	defer umap.mu.Unlock()

	for _, u := range umap.container {
		u.state = userStateBuffering
		u.isAnswered = false
	}
}

var (
	errUserNotFound = errors.New("user not found")
)
