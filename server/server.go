package server

import (
	"encoding/json"
	"github.com/flameous/anime-quiz/quiz"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	mu    sync.RWMutex
	rooms map[string]*quiz.Room
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*quiz.Room),
		mu:    sync.RWMutex{},
	}
}

func (s *Server) Start(addr string) error {

	// In production this URL serve by NGINX
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/quiz/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/quiz.html")
	})

	http.HandleFunc("/api/quiz_status", func(w http.ResponseWriter, r *http.Request) {
		if roomID, ok := r.URL.Query()["room_id"]; ok {
			room, _ := s.getRoomByID(roomID[0])

			js, err := json.Marshal(quiz.GetRoomStatus(room))

			if err != nil {
				log.Printf("server: quiz status: can't marshal response: %v", err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.Handle("/ws", s.handleWS())

	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleWS() http.Handler {
	h := func(conn *websocket.Conn) {
		// read for the first time
		var um quiz.UserMessage
		err := readMessage(conn, &um)
		if err != nil {
			return
		}

		if um.MessageType != quiz.UserMessageTypeHandShake {
			log.Println("server: ws: first msg is not a handshake")
			return
		}

		if um.RoomID == "" {
			log.Println("server: ws: roomID is empty")
			return
		}

		var userID string
		if um.Message == "" {
			log.Println("server: ws: message is empty")
			return
		}
		userID = um.Message

		room, ok := s.getRoomByID(um.RoomID)
		u := quiz.NewUser(userID, &quiz.WSConnection{Conn: conn})
		if ok {
			// is player (not admin)
			if room.RoomState != quiz.RoomStateInit {
				// todo: respond to user
				return
			}

			// add user to the room (sic!)
			room.AddUser(u)
		} else {
			// room is not created yet
			// user is admin
			room = s.addNewRoom(um.RoomID, userID)
			// add user to the room (sic!) code duplication (!!)
			room.AddUser(u)
			err = room.SendAdminNotify(userID) // notify admin
			if err != nil {
				log.Println(err)
			}
		}

		err = room.SendEnterNotifyToAll(userID)
		if err != nil {
			log.Printf("server: ws: can't notify all users %v", err)
		}

		for {
			var um quiz.UserMessage
			err = readMessage(conn, &um)
			if err != nil {
				return
			}
			room.HandleUserAction(u, um)
		}
	}

	return websocket.Handler(h)
}

func readMessage(conn *websocket.Conn, um *quiz.UserMessage) error {
	b := make([]byte, 2048)
	n, err := conn.Read(b)
	if err != nil {
		if err != io.EOF {
			log.Printf("server: ws: read user msg: %v", err)
		}
		return err
	}

	err = json.Unmarshal(b[:n], um)
	if err != nil {
		log.Printf("server: ws: msg unmarshal: %v, raw message: '%s'", err, b)
		return err
	}
	return nil
}

func (s *Server) getRoomByID(id string) (*quiz.Room, bool) {
	s.mu.RLock()
	r, ok := s.rooms[id]
	s.mu.RUnlock()

	return r, ok
}

func (s *Server) addNewRoom(roomID, adminID string) *quiz.Room {
	s.mu.Lock()
	room := quiz.NewRoom(roomID, adminID)
	s.rooms[roomID] = room
	s.mu.Unlock()

	return room
}
