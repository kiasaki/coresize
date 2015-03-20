package coresize

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kiasaki/batbelt/rest"
)

func (s *Server) handlePanic(w http.ResponseWriter, r *http.Request, err interface{}) {
	rest.SetInternalServerErrorResponse(w, err)
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	rest.SetNotFoundResponse(w)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rest.SetOKResponse(w, map[string]string{
		"app": "coresize",
	})
}

func (s *Server) handleFilePaths(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mappings := map[string]string{}

	for _, file := range s.Files {
		if s.Config.Hash {
			mappings[file.Name()] = file.NameWithHash()
		} else {
			mappings[file.Name()] = file.Name()
		}
	}

	rest.SetOKResponse(w, map[string]interface{}{
		"hashes": s.Config.Hash,
		"files":  mappings,
	})
}
