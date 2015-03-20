package coresize

import (
	"net/http"
	"strconv"

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

func (s *Server) handleImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	for _, file := range s.Files {
		var fileName string
		if s.Config.Hash {
			fileName = file.NameWithHash()
		} else {
			fileName = file.Name()
		}

		// If it's a match with requested file render else keep searching
		if fileName == ps.ByName("filename")[1:] {
			w.Header().Set("Content-Type", file.FileType())

			x, _ := strconv.Atoi(r.URL.Query().Get("x"))
			if x == 0 {
				x = 300
			}
			y, _ := strconv.Atoi(r.URL.Query().Get("y"))
			if y == 0 {
				y = 170
			}
			align := r.URL.Query().Get("align")
			if align == "" {
				align = "cc"
			}

			err := file.Render(w, x, y, align)
			if err != nil {
				rest.SetInternalServerErrorResponse(w, err)
			}
			return
		}
	}
	http.Error(w, "File not found", http.StatusNotFound)
}
