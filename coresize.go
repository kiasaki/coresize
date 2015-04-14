package coresize

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Config Config
	Router *httprouter.Router
	Cache  *Cache
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ParseFlags() {
	s.Config.ParseFlags()

	if s.Config.AwsSecretKey == "" || s.Config.AwsClientKey == "" {
		log.Fatal("AWS client and secret key are required configuration")
	}
	if s.Config.Bucket == "" {
		log.Fatal("AWS S3 Bucket name is a required configuration")
	}
}

func (s *Server) Setup() {
	// Routes
	router := httprouter.New()

	router.NotFound = s.handleNotFound
	router.PanicHandler = s.handlePanic

	router.GET("/", l(s.handleIndex))
	router.GET("/v1/i/*filename", l(s.handleImage))

	s.Router = router

	// File cache
	s.Cache = NewCache(s.Config)
	if err := s.Cache.Setup(); err != nil {
		log.Printf("Error setting up cache: %s\n", err.Error())
		os.Exit(1)
	}
}

func (s *Server) Run() {
	log.Printf("Listening on port %d\n", s.Config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.Config.Port), s.Router))
}

func (s *Server) SetupAndRun() {
	s.Setup()
	s.Run()
}

func l(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler(w, r, ps)
	}
}
