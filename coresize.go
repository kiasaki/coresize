package coresize

import (
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Config Config
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ParseFlags() {
	s.Config.ParseFlags()
}

func (s *Server) Run() {
	router := httprouter.New()

	router.NotFound = s.handleNotFound
	router.PanicHandler = s.handlePanic

	router.GET("/", s.handleIndex)

	log.Printf("Listening on port %d\n", s.Config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.Config.Port), router))
}
