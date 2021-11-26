package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Miyan0/bookings/pkg/config"
	"github.com/Miyan0/bookings/pkg/handlers"
	"github.com/Miyan0/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	//NOTE: change this to `true` in production
	app.InProduction = false

	// session initialization
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// loads all the templates (including layouts) in memory.
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	// stores the templates in memory to the cache
	app.TemplateCache = tc

	app.UseCache = false // development mode

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// now we need to render those templates
	render.NewTemplates(&app)

	// routes
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	// log to the console
	fmt.Println("Starting application on port", portNumber)

	// start the server
	// _ = http.ListenAndServe(portNumber, nil)

	// using third party router
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
