// insight project main.go
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

const DLS_FOLDER = "../dls"
const DLS_ZIP_FILE = "dls.zip"

// int mapping
const (
	IMAGE  int = 0
	FOLDER     = 1
)

type MainUI struct {
	FilesPath []FilePath
	Copyright string
}

type FilePath struct {
	P        string
	FileType int
}

func main() {

	fs := http.FileServer(http.Dir("./img"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", mainUI) // MUST BE THE LAST HandleFunc
	http.HandleFunc("/favicon.ico", favicon)
	http.ListenAndServe(":3333", nil)
}

func getURLEnd(r *http.Request) string {
	arr := strings.Split(r.URL.Path, "/")
	return arr[len(arr)-1]
}
func favicon(w http.ResponseWriter, r *http.Request) {
	//dat, _ := ioutil.ReadFile("favicon.ico")
	w.Write(FAVICON)
}
func mainUI(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./img")
	if err != nil {
		log.Fatal(err)
	}
	var filepaths []FilePath
	for _, file := range files {
		filepaths = append(filepaths, FilePath{P: file.Name(), FileType: IMAGE})
	}
	// Files are provided as a slice of strings.
	paths := []string{
		"main_ui.html",
	}

	t := template.Must(template.New("main_ui-tmpl").ParseFiles(paths...))
	err = t.ExecuteTemplate(w, "main_ui.html", MainUI{FilesPath: filepaths, Copyright: "© 2022 Sacha Liechti"})
	if err != nil {
		panic(err)
	}
}

func writeFileToResponseWriter(w http.ResponseWriter, filename string) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}
