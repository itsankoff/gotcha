package server

import (
    "log"
    "github.com/itsankoff/gotcha/common"
    "errors"
)


type Server struct {
    transports          map[string]Transport
    users               []*common.User
    connected           chan *common.User
    disconnected        chan *common.User
}

func New() *Server {
    return &Server{
        transports: make(map[string]Transport),
        users: make([]*common.User, 10),
        connected: make(chan *common.User),
        disconnected: make(chan *common.User),
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
    select {
    case user := <-s.connected:
        s.users = append(s.users, user)
        log.Println("Add user to server")
        go s.echoHandler(user)
    }
}

func (s *Server) userDisconnected() {
    log.Println("Start user disconnected observer")
    select {
    case user := <-s.disconnected:
        for i, u := range s.users {
            if u == user {
                s.users = append(s.users[:i], s.users[i+1:]...)
                log.Println("Remove user", user.Id)
                break
            }
        }
    }
}

func (s Server) echoHandler(user *common.User) {
    select {
    case msg := <-user.In:
        log.Println("Message from user", user.Id, msg.String())
        user.Out<-msg
    }
}

func (s *Server) Start(done <-chan interface{}) error {
    if len(s.transports) == 0 {
        return errors.New("Need to add transport before calling Start")
    }

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
