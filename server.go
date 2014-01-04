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
	"os/exec"
)

const srcDir string = "/Users/bks/src"

// Given a file list, return a list of mime types
func mimeType(filelist []string) []string {
	out := []string{}
	for _, fname := range filelist {
		str, err := exec.Command("file", "-b", "--mime-type", fname).Output()
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, strings.TrimSpace(string(str)))
	}
	return out
}

// Given an absolute path, generate a flat list of all files recursively
// (relative to the input directory)
func dirlist(absPath string) []string {
	out := []string{}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Couldn't find path ", path)
			return nil
		}

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(absPath, path)
			out = append(out, rel)
		}
		return nil
	}

	if err := filepath.Walk(absPath, walkFn); err != nil {
		log.Println("file walk failed", err)
	}
	return out
}

func writeJSON(slice interface{}, w http.ResponseWriter) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(&slice); err != nil {
		return err
	}
	return nil
}

type File map[string]string

type Repo struct {
	Name string `json:"name"`
	Readme string `json:"readme"`
	Files []File 	`json:"files,omitempty"`
}

func (r *Repo) SetFiles() {
	files := dirlist(filepath.Join(srcDir, r.Name))
	for i := 0; i < len(files); i++ {
		f := File{"name": files[i]}
		r.Files = append(r.Files, f)
	}
}

func findReadme(fios []os.FileInfo) string {
	for _, fi := range fios {
		path := fi.Name()
		if strings.HasPrefix(path, "README") || strings.HasPrefix(path, "readme") {
			return path
		}
	}
	return ""
}

func NewRepo(name string) *Repo {
	fios, err := ioutil.ReadDir(filepath.Join(srcDir, name))

	if err != nil {
		log.Println("Can't find repo ", name)
		return nil
	}

	return &Repo{
		Name: name,
		Readme: findReadme(fios)}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := r.URL.Path[1:]

	// /repo
	if strings.Index(path, "/") == -1 {

		r := NewRepo(path)
		r.SetFiles()
		if r == nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		if err := writeJSON(r, w); err != nil {
			log.Println(err)
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

	// /repo/plus/file.md
	} else {
		srcPath := filepath.Join(srcDir, path)
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

	fios, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Println(err)
		http.Error(w, "Oops", http.StatusInternalServerError)
		return
	}

	out := []*Repo{}
	for _, fi := range fios {
		if fi.IsDir() {
			r := NewRepo(fi.Name())
			out = append(out, r)
		}
	}
	writeJSON(out, w)
}

func main() {
	http.HandleFunc("/listing", listingHandler)
	http.HandleFunc("/", fileHandler)
	http.ListenAndServe(":8080", nil)
}
