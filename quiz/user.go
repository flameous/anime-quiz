package quiz

type User struct {
	id             string
	state          userState
	conn           connection
	answer         string
	isAnswerRight  bool
	arbitrageVotes map[string]bool
	arbitrageScore int
	score          int
}

func NewUser(userID string, c connection) *User {
	return &User{
		id:    userID,
		state: userStateJoined,
		conn:  c,
	}
}

func (u *User) sendMessageToUser(msg serverMessage) error {
	return u.conn.send(msg)
}

type userState int8

const (
	userStateJoined userState = iota // initial user state
	userStateBuffering
	userStateReadyToPlay // already buffered video OR playing
)
