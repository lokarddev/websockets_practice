package chat

import (
	"github.com/CloudyKit/jet/v6"
	"log"
	"net/http"
)

var views = jet.NewSet(jet.NewOSFileSystemLoader("./templates"), jet.InDevelopmentMode())

func Home(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, "home.html", nil); err != nil {
		log.Println(err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	if err = view.Execute(w, data, nil); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
