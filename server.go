package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"

	"github.com/gorilla/mux"
)

var notFoundErr = errors.New("Not Found")

func vueTemplate(filename string) *template.Template {
	name := path.Base(filename)
	return template.Must(template.New(name).Delims("[[", "]]").ParseFiles(filename))
}

func goTemplate(filename string) *template.Template {
	return template.Must(template.ParseFiles(filename))
}

func encodeKey(k *datastore.Key) string {
	return k.Encode()
}

func formatDate(d time.Time) string {
	return d.Format("01/02/2006")
}

var templateFns = map[string]interface{}{
	"encodeKey": encodeKey,
	"date":      formatDate,
}

func woTemplate(filename string) *template.Template {
	name := path.Base(filename)
	return template.Must(template.New(name).Funcs(templateFns).ParseFiles(filename))
}

var (
	indexTmpl       = goTemplate("tmpl/index.html")
	indexUserTmpl   = woTemplate("tmpl/indexUser.html")
	loginTmpl       = goTemplate("tmpl/login.html")
	newWorkoutTmpl  = vueTemplate("tmpl/newWorkout.html")
	showWorkoutTmpl = woTemplate("tmpl/showWorkout.html")
)

type indexPage struct {
	LoginURL string
}

type handler func(http.ResponseWriter, *http.Request) error

type userHandler func(http.ResponseWriter, *http.Request, *user.User, context.Context) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	c := appengine.NewContext(r)
	if err != nil && err == notFoundErr {
		log.Errorf(c, "Serving: %s: %s", r.URL.String(), err.Error())
		http.Error(w, "Not Found", 404)
	} else if err != nil {
		log.Errorf(c, "Serving: %s: %s", r.URL.String(), err.Error())
		http.Error(w, "Server Error", 500)
	}
}

func (h userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	u := user.Current(c)

	if u == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	err := h(w, r, u, c)

	if err != nil && err == notFoundErr {
		log.Errorf(c, "Serving: %s: %s", r.URL.String(), err.Error())
		http.Error(w, "Not Found", 404)
	} else if err != nil {
		log.Errorf(c, "Serving: %s: %s", r.URL.String(), err.Error())
		http.Error(w, "Server Error", 500)
	}
}

func render(w http.ResponseWriter, t *template.Template, data interface{}) error {
	buf := bytes.Buffer{}
	err := t.Execute(&buf, data)
	if err != nil {
		return err
	}

	if _, err := buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	u := user.Current(c)

	if u == nil {
		url, err := user.LoginURL(c, "/")
		if err != nil {
			return err
		}

		p := indexPage{
			LoginURL: url,
		}

		if err := render(w, indexTmpl, &p); err != nil {
			return err
		}

		return nil
	}

	//load up recent activities
	return indexUserHandler(w, r, u, c)

}

type userIndexPage struct {
	LogoutURL      string
	Name           string
	User           *User
	RecentWorkouts []*Workout
}

func indexUserHandler(w http.ResponseWriter, r *http.Request, gu *user.User, c context.Context) error {
	url, err := user.LogoutURL(c, "/")
	if err != nil {
		return err
	}

	u, err := FindOrCreateUser(c)
	if err != nil {
		return err
	}

	workouts, err := FindRecentWorkoutsForUser(c, u)
	if err != nil {
		return err
	}

	p := userIndexPage{
		LogoutURL:      url,
		Name:           gu.Email,
		User:           u,
		RecentWorkouts: workouts,
	}

	if err := render(w, indexUserTmpl, &p); err != nil {
		return err
	}

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	u := user.Current(c)

	if u == nil {
		if err := loginTmpl.Execute(w, nil); err != nil {
			return err
		}

		return nil
	}

	http.Redirect(w, r, "/", 302)
	return nil
}

type newWorkoutPage struct {
	Name      string
	LogoutURL string
}

func newWorkoutHandler(w http.ResponseWriter, r *http.Request, u *user.User, c context.Context) error {

	url, err := user.LogoutURL(c, "/")
	if err != nil {
		return err
	}

	p := newWorkoutPage{
		Name:      u.Email,
		LogoutURL: url,
	}

	if err := render(w, newWorkoutTmpl, &p); err != nil {
		return err
	}
	return nil
}

func createWorkoutHandler(w http.ResponseWriter, r *http.Request, gu *user.User, c context.Context) error {
	date := r.FormValue("date")
	t, err := time.Parse("01/02/2006", date)

	if err != nil {
		return newWorkoutHandler(w, r, gu, c)
	}

	wo := Workout{
		Type:    r.FormValue("type"),
		Details: r.FormValue("details"),
		Date:    t,
	}

	u, err := FindOrCreateUser(c)

	if err != nil {
		return err
	}

	if err := CreateWorkoutForUser(c, &wo, u); err != nil {
		return err
	}

	http.Redirect(w, r, fmt.Sprintf("/workouts/%s", wo.Key.Encode()), 302)
	return nil
}

type showWorkoutPage struct {
	Name      string
	LogoutURL string
	DateS     string
	Workout   *Workout
	User      *User
}

func showWorkoutHandler(w http.ResponseWriter, r *http.Request, gu *user.User, c context.Context) error {
	u, err := FindOrCreateUser(c)

	if err != nil {
		return err
	}

	vars := mux.Vars(r)
	woKey := vars["workoutKey"]

	wo, err := FindWorkoutForUser(c, u, woKey)
	if err != nil {
		return err
	}

	p := showWorkoutPage{
		User:    u,
		Workout: wo,
		DateS:   wo.Date.Format("01/02/2006"),
	}

	if err := render(w, showWorkoutTmpl, &p); err != nil {
		return err
	}
	return nil
}

func init() {
	r := mux.NewRouter()
	r.Handle("/", handler(indexHandler))
	r.Handle("/login", handler(loginHandler))
	r.Handle("/workouts/new", userHandler(newWorkoutHandler))
	r.Handle("/workouts/{workoutKey}", userHandler(showWorkoutHandler))
	r.Handle("/workouts", userHandler(createWorkoutHandler))
	http.Handle("/", r)
}
