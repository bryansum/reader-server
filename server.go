package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	// "fmt"
)

const srcDir string = "/Users/bks/src"

// Given an absolute path, generate a flat list of all files recursively
func dirlist(absPath string) ([]string, error) {
	out := []string{}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			fpath, _ := filepath.Rel(absPath, path)
			out = append(out, fpath)
		}
		return nil
	}

	if err := filepath.Walk(absPath, walkFn); err != nil {
		return nil, err
	}

	return out, nil
}

func writeJSONSlice(slice []string, w http.ResponseWriter) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(&slice); err != nil {
		return err
	}
	return nil
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := r.URL.Path[1:]
	srcPath := filepath.Join(srcDir, path)

	if strings.Index(path, "/") == -1 {
		slice, err := dirlist(srcPath)
		if err != nil {
			log.Println(err)
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
		if err := writeJSONSlice(slice, w); err != nil {
			log.Println(err)
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
	} else {
		if f, err := os.Open(srcPath); err != nil {
			log.Println(err)
			http.NotFound(w, r)
		} else {
			io.Copy(w, f)
		}
	}
}

func listingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if fileinfo, err := ioutil.ReadDir(srcDir); err != nil {
		log.Println(err)
		http.Error(w, "Oops", http.StatusInternalServerError)
	} else {
		out := []string{}
		for _, fi := range fileinfo {
			out = append(out, fi.Name())
		}
		writeJSONSlice(out, w)
	}
}

func main() {
	http.HandleFunc("/listing", listingHandler)
	http.HandleFunc("/", fileHandler)
	http.ListenAndServe(":8080", nil)
}
