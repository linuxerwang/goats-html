package main

import (
	"goats-html/examples/data"
	tmpl "goats-html/examples/shelf_view_html"
	"goats-html/goats/runtime"
	"log"
	"net/http"
)

var (
	shelf        = data.NewBookShelf()
	tmplSettings = &runtime.TemplateSettings{
		OmitDocType: false,
	}
)

func init() {
	// Init goats.
	goatsSettings := runtime.NewGoatsSettings()
	goatsSettings.PkgRoot = "."
	goatsSettings.TemplateDir = "goats-html/examples"
	runtime.InitGoats(goatsSettings)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	args := &tmpl.ShelfViewTemplateArgs{
		Shelf: shelf,
	}

	template := tmpl.NewShelfViewTemplate(w, tmplSettings)
	err := template.Render(args)
	if err != nil {
		log.Println("Failed to render template. ", err)
	}
}

func main() {
	http.HandleFunc("/", mainPageHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Println("Now you can visit http://localhost:8000/")
	}
}
