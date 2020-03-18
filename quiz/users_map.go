package quiz

import (
	"errors"
	"log"
	"math"
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

func (umap *usersMap) getScores() map[string]int {
	umap.mu.RLock()
	defer umap.mu.RUnlock()

	usersScore := make(map[string]int)
	for _, u := range umap.container {
		usersScore[u.id] = u.score
	}

	return usersScore
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

func (umap *usersMap) sendMessageForEachUserWithoutOne(userID string, msg serverMessage) error {
	umap.mu.RLock()
	defer umap.mu.RUnlock()

	var anyErr error

	for _, u := range umap.container {
		// Skip user
		if u.id == userID {
			continue
		}

		err := u.sendMessageToUser(msg)
		if err != nil {
			anyErr = err
		}
	}

	return anyErr
}

func (umap *usersMap) sendArbitrageApprovedToUsers() {
	umap.mu.RLock()
	defer umap.mu.RUnlock()

	for _, u := range umap.container {
		userCount := float64(len(umap.container))
		needVotes := int(math.Ceil(userCount / 2))

		if !u.isAnswerRight && u.arbitrageScore >= needVotes {
			u.score++

			srvMessage := serverMessage{
				MessageType: serverMessageTypeArbitrageApproved,
			}

			err := u.sendMessageToUser(srvMessage)
			if err != nil {
				log.Printf("room: send the arbitrage approve: %v", err)
			}
		}
	}
}

// on gameover
func (umap *usersMap) SendResultsForEachUser() error {
	usersScore := umap.getScores()

	msg := serverMessage{
		MessageType: serverMessageTypeGameOver,
		Message:     usersScore,
	}

	return umap.sendMessageForEachUser(msg)
}

func (umap *usersMap) setUsersToBuffered() {
	// reset all users
	umap.mu.Lock()
	defer umap.mu.Unlock()

	for _, u := range umap.container {
		u.state = userStateBuffering
		u.answer = ""
		u.isAnswerRight = false
		u.arbitrageScore = 0
		u.arbitrageVotes = make(map[string]bool)
	}
}

var (
	errUserNotFound = errors.New("user not found")
)
