package controllers

import (
	"net/http"

	"github.com/DhruvikDonga/tabsmooth/config/sessions"
	"github.com/DhruvikDonga/tabsmooth/views"
)

var index *views.View

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	session, err := sessions.Get(r)
	if err != nil {
		panic(err)
	}
	session.Save(r, w)
	//log.Println(session)
	index = views.NewView("bootstrap", "views/index.gohtml")
	index.Render(w, nil)

}
