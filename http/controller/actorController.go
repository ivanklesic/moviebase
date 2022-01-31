package http

import (
	"moviebase/moviebase/html"
	actorModel "moviebase/moviebase/model/entity/actor"
	movieModel "moviebase/moviebase/model/entity/movie"
	"moviebase/moviebase/service"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateActor(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		movies, err := movieModel.GetAll(0)
		if err != nil {
			panic(err)
		}
		tplData := make(html.TemplateData, 0)
		tplData["movies"] = movies
		html.Render(w, r, "actor/create", tplData)
		return
	} else if r.Method == "POST" {
		r.ParseMultipartForm(MB * 5)

		actor := actorModel.Actor{
			FirstName: strings.Trim(r.FormValue("fname"), " "),
			LastName:  strings.Trim(r.FormValue("lname"), " "),
		}

		ok, message, err := actor.Validate("create")

		if err != nil {
			panic(err)
		}

		if !ok {
			tplData := make(html.TemplateData, 0)
			tplData["message"] = message
			html.Render(w, r, "actor/create", tplData)
			return
		}

		movieIds := r.Form["movies[]"]

		if len(movieIds) > 0 {
			actor.MovieIds = movieIds

			for _, value := range movieIds {
				actor.Roles = append(actor.Roles, r.FormValue("role"+value))
			}
		}

		if _, ok := r.MultipartForm.File["img"]; ok {
			file, _, err := r.FormFile("img")

			if err != nil {
				panic(err)
			}

			filename, message, err := service.SaveImage(file, "actor")

			if err != nil {
				panic(err)
			}

			if message != "" {
				tplData := make(html.TemplateData, 0)
				tplData["message"] = message
				html.Render(w, r, "actor/create", tplData)
				return
			}
			actor.ImagePath = filename
		}

		actor.Create()

		http.Redirect(w, r, "/actor/list", http.StatusFound)
		return
	} else {
		panic("Method not allowed")
	}
}

func UpdateActor(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	actorId, err := strconv.Atoi(requestVars["id"])

	if err != nil {
		panic(err)
	}

	actor, err := actorModel.GetOne(actorId)

	if err != nil {
		panic(err)
	}

	if r.Method == "GET" {
		movies, err := movieModel.GetAll(0)
		if err != nil {
			panic(err)
		}
		tplData := make(html.TemplateData, 0)
		tplData["actor"] = actor
		tplData["movies"] = movies

		actor.LoadMovies()

		html.Render(w, r, "actor/create", tplData)
		return
	} else if r.Method == "POST" {
		r.ParseMultipartForm(MB * 5)

		actor := actorModel.Actor{
			Id:        actorId,
			FirstName: strings.Trim(r.FormValue("fname"), " "),
			LastName:  strings.Trim(r.FormValue("lname"), " "),
			ImagePath: actor.ImagePath,
		}

		ok, message, err := actor.Validate("update")

		if err != nil {
			panic(err)
		}

		if !ok {
			tplData := make(html.TemplateData, 0)
			tplData["message"] = message
			html.Render(w, r, "actor/create", tplData)
			return
		}

		movieIds := r.Form["movies[]"]

		if len(movieIds) > 0 {
			actor.MovieIds = movieIds

			for _, value := range movieIds {
				actor.Roles = append(actor.Roles, r.FormValue("role"+value))
			}
		}

		if _, ok := r.MultipartForm.File["img"]; ok {
			if actor.ImagePath != "" {
				err := os.Remove("public/images/actor/" + actor.ImagePath)
				if err != nil {
					panic(err)
				}
			}

			file, _, err := r.FormFile("img")
			if err != nil {
				panic(err)
			}
			filename, message, err := service.SaveImage(file, "actor")

			if err != nil {
				panic(err)
			}

			if message != "" {
				tplData := make(html.TemplateData, 0)
				tplData["message"] = message
				html.Render(w, r, "actor/create", tplData)
				return
			}

			actor.ImagePath = filename
		}

		actor.Update()

		http.Redirect(w, r, "/actor/list", http.StatusFound)
		return
	} else {
		panic("Method not allowed")
	}
}

func DeleteActor(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	actorId, err := strconv.Atoi(requestVars["id"])

	if err != nil {
		panic(err)
	}

	actor, err := actorModel.GetOne(actorId)

	if err != nil {
		panic(err)
	}

	actor.Delete()

	http.Redirect(w, r, "/actor/list", http.StatusFound)
}

func ListActors(w http.ResponseWriter, r *http.Request) {
	actors, err := actorModel.GetAll(0)
	if err != nil {
		panic(err)
	}
	if r.Method == "GET" {
		tplData := make(html.TemplateData, 0)
		tplData["actors"] = actors
		html.Render(w, r, "actor/list", tplData)
		return
	} else {
		panic("Method not allowed")
	}
}

func ActorDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		requestVars := mux.Vars(r)
		actorId, err := strconv.Atoi(requestVars["id"])

		if err != nil {
			panic(err)
		}

		actor, err := actorModel.GetOne(actorId)

		if err != nil {
			panic(err)
		}

		err = actor.LoadMovies()

		if err != nil {
			panic(err)
		}

		tplData := make(html.TemplateData, 0)
		tplData["actor"] = actor
		html.Render(w, r, "actor/details", tplData)
		return
	} else {
		panic("Method not allowed")
	}
}
