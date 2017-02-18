package server

import (
	"errors"
	"github.com/itsankoff/gotcha/common"
	"log"
)

type Server struct {
	transports      map[string]Transport
	users           []*common.User
	connected       chan *common.User
	disconnected    chan *common.User
	aggregate       chan *common.Message
	messageHandlers map[string]chan *common.Message

	outputStore  *OutputStore
	history      *History
	contactStore *ContactStore
	authRegistry *AuthRegistry
	fileStore    *FileStore

	control   *Control
	messanger *Messanger
}

func New(config *Config) *Server {
	s := &Server{
		transports:      make(map[string]Transport),
		users:           make([]*common.User, 10),
		connected:       make(chan *common.User),
		disconnected:    make(chan *common.User),
		aggregate:       make(chan *common.Message),
		messageHandlers: make(map[string]chan *common.Message),
	}

	s.contactStore = NewContactStore()
	s.authRegistry = NewAuthRegistry()
	s.outputStore = NewOutputStore()
	s.fileStore = NewFileStore(config.FileServerFolder,
		config.FileServerHost, config.FileServerPath)

	historyInput := make(chan *common.Message)
	s.history = NewHistory(historyInput, s.outputStore)

	// register control handler
	controlInput := make(chan *common.Message)
	s.control = NewControl(controlInput, s.outputStore,
		s.contactStore, s.authRegistry)
	s.messageHandlers["control"] = controlInput

	// register message handler
	messangerInput := make(chan *common.Message)

	s.messanger = NewMessanger(messangerInput, s.history,
		s.outputStore, s.fileStore)
	s.messageHandlers["message"] = messangerInput
	s.messageHandlers["file"] = messangerInput
	s.messageHandlers["history"] = historyInput

	return s
}

func (s *Server) startRouter() {
	for {
		select {
		case msg := <-s.aggregate:
			cmdType := msg.CmdType()
			handler, ok := s.messageHandlers[cmdType]
			if ok {
				handler <- msg
			} else {
				log.Println("No registered handler for cmd type", cmdType)
			}
		}
	}
}

func (s *Server) aggregateMessages(user *common.User) {
	// 0: connected
	// 1: registered
	// 2: authenticated
	var state int
	for {
		select {
		case msg := <-user.In:
			if msg == nil {
				log.Println("Nil message in aggregate for id %s userId %s",
					user.Id, user.UserId)
				s.outputStore.RemoveOutput(user.UserId)
				return
			}

			switch state {
			case 0:
				if msg.CmdType() != "auth" && msg.Cmd() != "register" {
					log.Printf("Wrong message %s for state %d", msg.Cmd(), state)
					user.Disconnect()
					return
				}

				packet, err := msg.ParseJsonData()
				if err != nil {
					log.Println("Failed to parse register message data", err)
					user.Disconnect()
					return
				}

				username := packet["username"].(string)
				pass := packet["password"].(string)
				userId, registered := s.authRegistry.Register(username, pass)
				if !registered {
					log.Println("Registration for %s failed", username)
					user.Disconnect()
					return
				}

				response := common.NewResponse(msg, "user_id", userId)
				user.Out <- response
				state = 1
			case 1:
				if msg.CmdType() != "auth" || msg.Cmd() != "auth" {
					log.Printf("Wrong message %s for state %d", msg.Cmd(), state)
					user.Disconnect()
					return
				}

				packet, err := msg.ParseJsonData()
				if err != nil {
					log.Println("Failed to parse auth message data", err)
					user.Disconnect()
					return
				}

				userId := packet["user_id"].(string)
				pass := packet["password"].(string)
				authenticated := s.authRegistry.Authenticate(userId, pass)
				if !authenticated {
					log.Println("Authentication for %s failed", userId)
					user.Disconnect()
					return
				}

				user.UserId = userId
				s.outputStore.AddOutput(user.UserId, user.Out)
				state = 2

				response := common.NewResponse(msg, "authenticated", true)
				user.Out <- response
			case 2:
				log.Println("Forward message to aggregate", msg.From())
				s.aggregate <- msg
			}
		}
	}
}

func (s *Server) AddTransport(host string, t Transport) error {
	if host == "" {
		return errors.New("Can't add transport for an empty host")
	}

	if t == nil {
		return errors.New("Can't add nil transport")
	}

	_, ok := s.transports[host]
	if ok {
		// prevent adding multiple transports for the same url
		return errors.New("Try to add multiple transports for " + host)
	}

	s.transports[host] = t
	log.Println("Add transport for", host)
	return nil
}

func (s *Server) RemoveTransport(host string) error {
	_, ok := s.transports[host]
	if ok {
		delete(s.transports, host)
		log.Println("Remove transport for", host)
		return nil
	}

	return errors.New("No trasport for host " + host)
}

func (s *Server) userConnected() {
	log.Println("Start user connected observer")
	for {
		select {
		case user := <-s.connected:
			s.users = append(s.users, user)
			log.Println("Add user to server")
			go s.aggregateMessages(user)
		}
	}
}

func (s *Server) userDisconnected() {
	log.Println("Start user disconnected observer")
	for {
		select {
		case user := <-s.disconnected:
			for i, u := range s.users {
				if u == user {
					s.users = append(s.users[:i], s.users[i+1:]...)
					log.Println("Remove user from server", user.Id)
					break
				}
			}
		}
	}
}

func (s Server) echoHandler(user *common.User) {
	select {
	case msg := <-user.In:
		user.Out <- msg
	}
}

func (s *Server) Start(done <-chan interface{}) error {
	if len(s.transports) == 0 {
		return errors.New("Need to add transport before calling Start")
	}

	go s.startRouter()
	go s.userConnected()
	go s.userDisconnected()
	for url, t := range s.transports {
		log.Println("Start transport for", url)
		t.OnUserConnected(s.connected)
		t.OnUserDisconnected(s.disconnected)
		go t.Start(url, done)
	}

	<-done
	close(s.connected)
	close(s.disconnected)
	close(s.aggregate)
	return nil
}

func (s *Server) StartAsync() (chan interface{}, error) {
	if len(s.transports) == 0 {
		return nil, errors.New("Need to add transport before calling Start")
	}

	done := make(chan interface{})
	go s.Start(done)
	return done, nil
}
