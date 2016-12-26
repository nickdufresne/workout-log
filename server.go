package main

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

var notFoundErr = errors.New("Not Found")

var (
	indexTmpl      = template.Must(template.ParseFiles("tmpl/index.html"))
	indexUserTmpl  = template.Must(template.ParseFiles("tmpl/indexUser.html"))
	loginTmpl      = template.Must(template.ParseFiles("tmpl/login.html"))
	newWorkoutTmpl = template.Must(template.ParseFiles("tmpl/newWorkout.html"))
)

type indexPage struct {
	LoginURL string
}

type userIndexPage struct {
	LogoutURL string
	Name      string
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

	url, err := user.LogoutURL(c, "/")
	if err != nil {
		return err
	}

	p := userIndexPage{
		LogoutURL: url,
		Name:      u.Email,
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

func createWorkoutHandler(w http.ResponseWriter, r *http.Request, u *user.User, c context.Context) error {
	w.Write([]byte(`hi`))
	return nil
}

func init() {
	http.Handle("/", handler(indexHandler))
	http.Handle("/login", handler(loginHandler))
	http.Handle("/workouts/new", userHandler(newWorkoutHandler))
	http.Handle("/workouts", userHandler(createWorkoutHandler))
}
