package coresize

import (
	"log"
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
	// return supported endpoints with url templates
	rest.SetOKResponse(w, map[string]interface{}{
		"v1": map[string]string{
			"root_url":  "/",
			"image_url": "/v1/image/{file_name}?hash={file_hash}&height={height}&width={width}&allign={allignment}",
		},
	})
}

func (s *Server) handleImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestFilename := ps.ByName("filename")[1:]

	if cacheFile, ok := s.Cache.Get(requestFilename); ok {
		// extract url params
		height, _ := strconv.Atoi(r.URL.Query().Get("height"))
		width, _ := strconv.Atoi(r.URL.Query().Get("width"))
		align := r.URL.Query().Get("align")
		if align == "" {
			align = "cc"
		}

		// fetch and render, or, read from disk
		err := cacheFile.Render(w, width, height, align)
		if err != nil {
			log.Printf("Error rendering %s: %s\n", requestFilename, err.Error())
			rest.SetInternalServerErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", cacheFile.FileType())
		return
	}
	http.Error(w, "File not found", http.StatusNotFound)
}
