package server

import (
    "log"
    "encoding/json"
    "errors"
    "github.com/itsankoff/gotcha/common"
)


type Server struct {
    transports          map[string]Transport
    users               []*common.User
    connected           chan *common.User
    disconnected        chan *common.User
    aggregate           chan *common.Message
    messageHandlers     map[string]chan *common.Message

    control             *Control
    messanger           *Messanger
    history             *History
    authRegistry        *AuthRegistry
    outputStore         *OutputStore
}

func New() *Server {
    s := &Server{
        transports: make(map[string]Transport),
        users: make([]*common.User, 10),
        connected: make(chan *common.User),
        disconnected: make(chan *common.User),
        aggregate: make(chan *common.Message),
        messageHandlers: make(map[string]chan *common.Message),
    }

    s.outputStore = NewOutputStore()
    s.history = NewHistory()

    // register control handler
    controlInput := make(chan *common.Message)
    s.control = NewControl(controlInput, s.outputStore)
    s.messageHandlers["control"] = controlInput

    // register message handler
    messangerInput := make(chan *common.Message)

    s.messanger = NewMessanger(messangerInput, s.history, s.outputStore)
    s.messageHandlers["message"] = messangerInput
    s.messageHandlers["file"] = messangerInput

    s.authRegistry = NewAuthRegistry()

    return s
}

func (s *Server) startRouter() {
    for {
        select {
        case msg := <-s.aggregate:
            log.Println("Message in aggregate", msg.From())
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
            switch state {
            case 0:
                if msg.CmdType() != "auth" && msg.Cmd() != "register" {
                    log.Printf("Wrong message %s for state %d", msg.Cmd(), state)
                    user.Disconnect()
                    return
                }

                var packet map[string]interface{}
                err := json.Unmarshal([]byte(msg.String()), &packet)
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

                response := common.NewResponse(msg, userId)
                user.Out <- response
                state = 1
            case 1:
                if msg.CmdType() != "auth" || msg.Cmd() != "auth" {
                    log.Printf("Wrong message %s for state %d", msg.Cmd(), state)
                    user.Disconnect()
                    return
                }

                var packet map[string]interface{}
                err := json.Unmarshal([]byte(msg.String()), &packet)
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

                s.outputStore.AddOutput(userId, user.Out)
                state = 2

                response := common.NewResponse(msg, "authenticated")
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
            s.outputStore.AddOutput(user.Id, user.Out)
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
                    s.outputStore.RemoveOutput(user.Id)
                    log.Println("Remove user", user.Id)
                    break
                }
            }
        }
    }
}

func (s Server) echoHandler(user *common.User) {
    select {
    case msg := <-user.In:
        user.Out<-msg
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
