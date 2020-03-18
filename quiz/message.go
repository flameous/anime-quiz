package quiz

type serverMessage struct {
	ID          string      `json:"id,omitempty"` // for debug
	MessageType string      `json:"message_type"`
	Message     interface{} `json:"message"`
}

type UserMessage struct {
	ID          string `json:"id,omitempty"`
	RoomID      string `json:"room_id,omitempty"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

const (
	UserMessageTypeHandShake = "userHandShake"
	UserMessageTypeNotify    = "userNotify"
	UserMessageTypeAnswer    = "userAnswer"
	UserMessageTypeArbitrage = "userArbitrage"

	serverMessageTypeAdminNotify       = "serverAdminNotify"
	serverMessageTypeEnterNotify       = "serverEnterNotify"
	serverMessageTypeStartGame         = "serverStartGame"
	serverMessageTypeSendVideo         = "serverSendVideo"
	serverMessageTypeStartPlaying      = "serverStartPlaying"
	serverMessageTypeAnswer            = "serverAnswer"
	serverMessageTypeArbitrage         = "serverArbitrage"
	serverMessageTypeArbitrageApproved = "serverArbitrageResult"
	serverMessageTypeGameOver          = "serverGameOver"
)
