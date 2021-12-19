package main

import (
	"context"
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DhruvikDonga/tabsmooth/config/migrations"
	"github.com/DhruvikDonga/tabsmooth/config/sessions"
	"github.com/DhruvikDonga/tabsmooth/controllers"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	gob.Register(sessions.Flash{})
	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:5000"
	}
	dns := os.Getenv("DNS")
	if dns != "" {
		migrations.DNS = dns
	}
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	migrations.IntialMigration()

	router := mux.NewRouter()
	CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))

	//---------------Routes--------------
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.PathPrefix("/storage/").Handler(http.FileServer(http.Dir("./storage")))

	router.HandleFunc("/", controllers.Index).Methods("GET")
	router.HandleFunc("/blogs", controllers.GetBlogs).Methods("GET")
	router.HandleFunc("/createblog", controllers.CreatePost).Methods("GET")
	router.HandleFunc("/blog/{slug}", controllers.SingleBlog).Methods("GET")
	router.HandleFunc("/postblog", controllers.PostBlog).Methods("POST")
	router.HandleFunc("/your-blogs", controllers.GetUserBlogs).Methods("GET")
	router.HandleFunc("/edit/{slug}", controllers.EditBlogForm).Methods("GET")
	router.HandleFunc("/posteditblog", controllers.EditBlog).Methods("POST")

	router.HandleFunc("/register", controllers.RegisterPage).Methods("GET")
	router.HandleFunc("/postregister", controllers.PostRegister).Methods("POST")
	router.HandleFunc("/login", controllers.LoginPage).Methods("GET")
	router.HandleFunc("/postlogin", controllers.PostLogin).Methods("POST")
	router.HandleFunc("/logout", controllers.PostLogout).Methods("POST")
	router.HandleFunc("/profile", controllers.GetUserProfile).Methods("GET")
	router.HandleFunc("/resetprofile", controllers.ResetProfile).Methods("POST")
	//-----------------EndRoutes------------

	srv := &http.Server{
		Addr:         port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      CSRF(router),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)

}
