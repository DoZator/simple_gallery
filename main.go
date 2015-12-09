package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	images      = []Image{}
	extensions  = []string{".jpg", ".jpeg"}
	thumbs_path = flag.String("t", "thumbs", "Thumbs path")
	img_path    = flag.String("i", "images", "Images path")
	port        = flag.String("p", "9000", "Server port")
)

func main() {

	flag.Parse()

	images_path, err := filepath.Abs(*img_path)
	if err != nil {
		panic(err)
	}

	prepareImagesForPath(images_path)

	if len(images) > 0 {
		generateThumbs()
	}

	http.Handle("/o/", http.StripPrefix(path.Join("/o", *img_path), http.FileServer(http.Dir(*img_path))))
	http.Handle("/t/", http.StripPrefix(path.Join("/t", *thumbs_path), http.FileServer(http.Dir(*thumbs_path))))
	http.HandleFunc("/favicon.ico", http.NotFound)

	http.HandleFunc("/", handler)

	fmt.Println(fmt.Sprintf("Starting server on %s", *port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

func prepareImagesForPath(path string) {
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isValidImageForFilePath(filePath) {
			baseName := filepath.Base(filePath)
			name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
			images = append(images, Image{Path: filePath, Name: name})
		}
		return nil
	})
}

func isValidImageForFilePath(filePath string) bool {
	if strings.HasPrefix(filepath.Base(filePath), ".") {
		return false
	}

	for _, e := range extensions {
		if e == strings.ToLower(path.Ext(filePath)) {
			return true
		}
	}

	return false
}

func generateThumbs() {
	if _, err := os.Stat(*thumbs_path); os.IsNotExist(err) {
		os.Mkdir(*thumbs_path, 0777)
	}

	path, err := filepath.Abs(*thumbs_path)
	if err != nil {
		return
	}

	clearThumbsDir(path)

	for i := range images {
		images[i].GenerateThumb(path)
	}
}

func clearThumbsDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("foo").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
</head>
<body>
{{ range . }}<a title="{{ .Name }}" href="/o{{ .Path }}">
<img alt="" width=230 height=230 src="/t/thumbs/{{ .ThumbName }}">
</a>
{{ end }}
</body>
</html>`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, images)
}
