package server

import (
    "log"
    "github.com/itsankoff/gotcha/util"
    "errors"
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

func (s *Server) userConnected(user *util.User) {
    s.users = append(s.users, user)
    // TODO: Generate unique user id
    log.Println("Add user")
}

func (s *Server) userDisconnected(user *util.User) {
    for i, u := range s.users {
        if u == user {
            s.users = append(s.users[:i], s.users[i+1:]...)
            log.Println("Remove user", user.Id)
            break
        }
    }
}

func (s *Server) Start(done <-chan interface{}) error {
    if len(s.transports) == 0 {
        return errors.New("Need to add transport before calling Start")
    }

    for url, t := range s.transports {
        log.Println("Start transport for", url)
        t.OnUserConnected(s.userConnected)
        t.OnUserDisconnected(s.userDisconnected)
        go t.Start(url, done)
    }

    <-done
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
