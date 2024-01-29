package main

import (
	"encoding/json"
	"net/http"
	"os"
	"word-search-in-files/pkg/searcher"
)

func wordSearch(w http.ResponseWriter, r *http.Request) {
	var errJson error
	var response []byte

	params := r.URL.Query()
	name := params.Get("word")

	fsys := os.DirFS("./examples")

	s := &searcher.Searcher{
		FS: fsys,
	}

	gotFiles, err := s.Search(name)

	if err != nil {
		response, errJson = json.Marshal(map[string]string{"error": err.Error()})
	} else {
		response, errJson = json.Marshal(gotFiles)
	}

	if errJson != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	http.HandleFunc("/files/search", wordSearch)
	http.ListenAndServe(":8080", nil)
}
