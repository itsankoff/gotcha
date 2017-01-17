package gotcha

import (
    "fmt"
)

type Server struct {

}

func New() *Server {
    return &Server{}
}

func (s *Server) Start(host string) {
    fmt.Println("Hello server")
}
