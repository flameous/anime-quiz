package quiz

type serverMessage struct {
	ID          string `json:"id,omitempty"` // for debug
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
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

	serverMessageTypeAdminNotify  = "serverAdminNotify"
	serverMessageTypeEnterNotify  = "serverEnterNotify"
	serverMessageTypeSendVideo    = "serverSendVideo"
	serverMessageTypeStartPlaying = "serverStartPlaying"
	serverMessageTypeAnswer       = "serverAnswer"
	serverMessageTypeGameOver     = "serverGameOver"
)
