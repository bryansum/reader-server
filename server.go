package main

import (
	"mime"
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

// True if ignored
func gitIgnore(fname string) bool {
	// 0 means ignored; 1 means none ignored
	return exec.Command("git", "check-ignore", fname).Run() == nil;
}

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

func dirlist(repo string) []string {
	gitDir := filepath.Join(repo, ".git")
	out, err := exec.Command("git", "--git-dir", gitDir, "ls-files").Output()
	if err != nil {
		log.Println("Couldn't find repo", repo)
	}
	return strings.Fields(string(out))
}

func writeJSON(slice interface{}, w http.ResponseWriter) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(&slice); err != nil {
		return err
	}
	return nil
}

type File struct {
	Name string `json:"name"`
	mimeType string `json:"mime-type"`
	Size int64 `json:"size"`
}

func NewFile(path string) File {
	// split repo from file path
	name := path[strings.Index(path, "/")+1:]
	var size int64
	if fi, err := os.Stat(path); err == nil {
		size = fi.Size()
	}
	mime := mime.TypeByExtension(filepath.Ext(path))
	return File{Name: name, mimeType: mime, Size: size}
}

type Repo struct {
	Name string `json:"name"`
	Readme File `json:"readme"`
	Files []File 	`json:"files,omitempty"`
}

func (r *Repo) SetFiles() {
	files := dirlist(r.Name)
	for _, fname := range files {
		r.Files = append(r.Files, NewFile(filepath.Join(r.Name, fname)))
	}
}

func findReadme(fios []os.FileInfo) File {
	for _, fi := range fios {
		name := fi.Name()
		if strings.HasPrefix(name, "README") || strings.HasPrefix(name, "readme") {
			return File{Name: name}
		}
	}
	return File{}
}

func NewRepo(name string) *Repo {
	fios, err := ioutil.ReadDir(name)

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
		if f, err := os.Open(path); err != nil {
			log.Println(err)
			http.NotFound(w, r)
		} else {
			mimeType := NewFile(path).mimeType
			w.Header().Set("Content-Type", mimeType)
			io.Copy(w, f)
		}
	}
}

func listingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fios, err := ioutil.ReadDir(".")
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

func addExtensions() {
	mime.AddExtensionType(".c", "text/x-csrc")
	mime.AddExtensionType(".h", "text/x-csrc")
	mime.AddExtensionType(".cpp", "text/x-c++src")
	mime.AddExtensionType(".hpp", "text/x-c++src")

	mime.AddExtensionType(".java", "text/x-java")

	mime.AddExtensionType(".m", "text/x-csrc")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".json", "application/json")

	mime.AddExtensionType(".sh", "text/x-sh")

	mime.AddExtensionType(".go", "text/x-go")

	mime.AddExtensionType(".php", "text/x-php")
	mime.AddExtensionType(".py", "text/x-python")
	mime.AddExtensionType(".rb", "text/x-ruby")

	mime.AddExtensionType(".md", "text/x-markdown")
	mime.AddExtensionType(".markdown", "text/x-markdown")

	mime.AddExtensionType(".coffee", "text/x-coffeescript")

	mime.AddExtensionType(".scss", "text/x-scss")
	mime.AddExtensionType(".less", "text/x-less")

	mime.AddExtensionType(".erb", "application/x-erb")
	mime.AddExtensionType(".ejs", "application/x-ejs")
}

func main() {
	addExtensions()

	if os.Chdir(srcDir) != nil {
		log.Println("Failed changing directory to", srcDir)
	}

	http.HandleFunc("/listing", listingHandler)
	http.HandleFunc("/", fileHandler)
	http.ListenAndServe(":8080", nil)
}
