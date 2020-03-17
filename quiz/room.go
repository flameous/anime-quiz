package quiz

import (
	"errors"
	"log"
	"time"
)

type roomState int8

const (
	RoomStateInit roomState = iota
	RoomStateReadyToPlay
	RoomStatePlayVideo
	RoomStatePlayShowAnswer
	RoomStateFinished
)

type Room struct {
	id            string
	users         *usersMap
	adminID       string
	RoomState     roomState
	currentQuizID int
	allQuizzes    []quiz
}

func NewRoom(roomID, adminID string) *Room {
	return &Room{
		id:            roomID,
		users:         newUsersMap(),
		adminID:       adminID,
		RoomState:     RoomStateInit,
		currentQuizID: 0,
		allQuizzes:    getShuffledQuizzes(hardcodedQuizzes),
	}
}

func (r *Room) AddUser(u *User) {
	r.users.addUser(u)
}

func (r *Room) HandleUserAction(u *User, message UserMessage) {

	switch r.RoomState {
	case RoomStateInit:
		r.handleRoomInit(u, message)

	case RoomStateReadyToPlay:
		r.handleRoomReady(u, message)

	case RoomStatePlayVideo:
		r.handleRoomPlaying(u, message)

	case RoomStatePlayShowAnswer:
		r.handleRoomShowAnswer(u, message)

	case RoomStateFinished:
		// gtfo
		r.handleFinishedRoom(u, message)
	default:
		// gtfo for real
		log.Printf("room: handle user actions: undefined room state: %v", r.RoomState)
	}
}

func (r *Room) handleRoomInit(u *User, msg UserMessage) {
	if u.id != r.adminID {
		// if room is ready to play
		// ignore other users input
		return
	}
	if msg.MessageType == UserMessageTypeNotify &&
		msg.Message == "startGame" {
		r.RoomState = RoomStateReadyToPlay

		err := r.sendCurrentVideoToAllUsers()
		if err != nil {
			log.Printf("room handler: handleRoomInit: %v", err)
		}
	}
}

// users are waiting for loading and buffering video
func (r *Room) handleRoomReady(u *User, msg UserMessage) {
	if msg.MessageType != UserMessageTypeNotify {
		return
	}

	u.state = userStateReadyToPlay // todo: refactor: make setter func
	if r.users.isAllReady() {
		r.RoomState = RoomStatePlayVideo

		// 1. send command to start playing video
		srvMessage := serverMessage{
			MessageType: serverMessageTypeStartPlaying,
		}
		err := r.users.sendMessageForEachUser(srvMessage)

		if err != nil {
			log.Printf("room handler: handleRoomReady: err")
		}

		// 2. run goroutine with timeout
		// after N secs send right answer and stop receiving answers
		go r.startQuiz()
	}
}

func (r *Room) handleRoomPlaying(u *User, msg UserMessage) {
	// receive answers
	// but just once
	if msg.MessageType != UserMessageTypeAnswer || len(u.answer) > 0 {
		return
	}

	u.answer = msg.Message
}

func (r *Room) handleRoomShowAnswer(u *User, msg UserMessage) {
	// Arbitrage
	if msg.MessageType != UserMessageTypeArbitrage {
		return
	}

	userID := msg.Message

	// go away if you already vote
	if _, ok := u.arbitrageVotes[userID]; ok {
		return
	}

	r.users.getUser(userID).arbitrageScore++
	u.arbitrageVotes[userID] = true
}

func (r *Room) handleFinishedRoom(u *User, msg UserMessage) {
	// what do you want here?
}

func (r *Room) SendAdminNotify(userID string) error {
	msg := serverMessage{
		MessageType: serverMessageTypeAdminNotify,
		Message:     "ti glavnuy!",
	}
	return r.users.sendMessageToUser(userID, msg)
}

func (r *Room) SendEnterNotifyToAll(userID string) error {
	r.users.mu.RLock()
	defer r.users.mu.RUnlock()

	msg := serverMessage{
		MessageType: serverMessageTypeEnterNotify,
		Message: map[string]interface{}{
			"count":   len(r.users.container),
			"user_id": userID,
		},
	}

	return r.users.sendMessageForEachUser(msg)
}

func (r *Room) sendCurrentVideoToAllUsers() error {
	if r.currentQuizID > len(r.allQuizzes)-1 {
		return errors.New("video index out of bounds")
	}

	currentQuiz := r.allQuizzes[r.currentQuizID]

	msg := serverMessage{
		MessageType: serverMessageTypeSendVideo,
		Message: map[string]interface{}{
			"video_id": currentQuiz.videoSource,
			"start":    currentQuiz.start,
		},
	}

	r.users.setUsersToBuffered()
	return r.users.sendMessageForEachUser(msg)
}

const (
	showVideoDuration  = time.Second * 10
	showAnswerDuration = time.Second * 5
)

func (r *Room) startQuiz() {
loop:
	for {
		switch r.RoomState {
		case RoomStatePlayVideo:
			time.Sleep(showVideoDuration)
			r.RoomState = RoomStatePlayShowAnswer

			q := r.allQuizzes[r.currentQuizID]

			srvMessage := serverMessage{
				MessageType: serverMessageTypeAnswer,
				Message:     q.title,
			}

			err := r.users.sendMessageForEachUser(srvMessage)
			if err != nil {
				log.Printf("room: send answer: %v", err)
			}

			// Check answers and send arbitrage
			for _, u := range r.users.container {
				if q.isAnswerRight(u.answer) {
					u.score++
					u.isAnswerRight = true
				} else {
					if len(u.answer) > 0 {
						srvMessage = serverMessage{
							MessageType: serverMessageTypeArbitrage,
							Message: map[string]interface{}{
								"user_id": u.id,
								"answer":  u.answer,
							},
						}

						err = r.users.sendMessageForEachUserWithoutOne(u.id, srvMessage)
						if err != nil {
							log.Printf("room: send an arbitrage: %v", err)
						}
					}
				}
			}

		case RoomStatePlayShowAnswer:
			time.Sleep(showAnswerDuration)
			if r.currentQuizID >= len(r.allQuizzes)-1 { // this quiz is the last one
				r.RoomState = RoomStateFinished
				r.users.sendArbitrageApprovedToUsers()

				err := r.users.SendResultsForEachUser()
				if err != nil {
					log.Printf("room: send game over: %v", err)
				}

			} else {
				// we have more quizzes to run
				r.RoomState = RoomStateReadyToPlay
				r.users.sendArbitrageApprovedToUsers()

				r.currentQuizID++
				err := r.sendCurrentVideoToAllUsers()
				if err != nil {
					log.Printf("room: send new video: %v", err)
				}
			}

		default:
			break loop
		}
	}
}
