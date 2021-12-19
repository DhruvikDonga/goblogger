package sessions

import (
	"encoding/gob"
	"net/http"

	"github.com/DhruvikDonga/tabsmooth/models"
	"github.com/gorilla/sessions"
	gsessions "github.com/gorilla/sessions"
)

var store = gsessions.NewCookieStore([]byte("contentplay-secret-key"))

type Flash struct {
	Key   string
	Value string
}

func init() {

	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   3600 * 8,
		Secure:   false,
		HttpOnly: true,
		SameSite: 0,
	}
	gob.Register(models.Auth{})

}
func Get(req *http.Request) (*gsessions.Session, error) {
	return store.Get(req, "newsession5")
}

func GetNamed(req *http.Request, name string) (*gsessions.Session, error) {
	return store.Get(req, name)
}
