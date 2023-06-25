package spa

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Handler struct {
	StaticPath string
	IndexPath  string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	dir, file := path.Split(r.URL.Path)
	reqPath := dir + file
	// log.Println(reqPath)
	if !path.IsAbs(reqPath) {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, "404 sorry not found.", http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	reqPath = filepath.Join(h.StaticPath, reqPath)

	// check whether a file exists at the given path
	_, err := os.Stat(reqPath)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop the execution
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}
