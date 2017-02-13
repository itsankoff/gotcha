package server

import (
    "log"
    "github.com/itsankoff/gotcha/util"
)


type Server struct {
    transports  map[string]Transport
    users       []*util.User
}

func New() *Server {
    return &Server{
        transports: make(map[string]Transport),
        users: make([]*util.User, 10),
    }
}

func (s *Server) AddTransport(host string, t Transport) {
    _, ok := s.transports[host]
    if ok {
        // prevent adding multiple transports for the same url
        log.Println("Try to add multiple transports for same host")
        return
    }

    s.transports[host] = t
    log.Println("Add transport", s.transports)
}

func (s *Server) RemoveTransport(host string, transport Transport) {
    _, ok := s.transports[host]
    if ok {
        delete(s.transports, host)
        log.Println("Remove transport for", host)
    }
}

func (s *Server) userConnected(user *util.User) {
    s.users = append(s.users, user)
    // TODO: Generate unique user id
    log.Println("Add user")
}

func (s *Server) userDisconnected(user *util.User) {
    for i, u := range s.users {
        if u == user {
            s.users = append(s.users[:i], s.users[i+1:]...)
            log.Println("Remove user", user.UserId)
            break
        }
    }
}

func (s *Server) Start(done <-chan interface{}) {
    for url, t := range s.transports {
        log.Println("Start transport for", url)
        t.OnUserConnected(s.userConnected)
        t.OnUserDisconnected(s.userDisconnected)
        go t.Start(url, done)
    }

    log.Println("Wait to close done channel")
    <-done
}

func (s *Server) StartAsync() (done <-chan interface{}) {
    go s.Start(done)
    return done
}
