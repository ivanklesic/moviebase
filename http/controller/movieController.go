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

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		actors, err := actorModel.GetAll(0)
		if err != nil {
			panic(err)
		}
		tplData := make(html.TemplateData, 0)
		tplData["actors"] = actors
		html.Render(w, r, "movie/create", tplData)
		return
	} else if r.Method == "POST" {
		r.ParseMultipartForm(MB * 5)

		movie := movieModel.Movie{
			Name:        strings.Trim(r.FormValue("name"), " "),
			Description: strings.Trim(r.FormValue("description"), " "),
			Year: strings.Trim(r.FormValue("year"), " "),
		}

		ok, message, err := movie.Validate("create")

		if err != nil {
			panic(err)
		}

		if !ok {
			tplData := make(html.TemplateData, 0)
			tplData["message"] = message
			html.Render(w, r, "movie/create", tplData)
			return
		}

		actorIds := r.Form["actors[]"]

		if len(actorIds) > 0 {
			movie.ActorIds = actorIds

			for _, value := range actorIds {
				movie.Roles = append(movie.Roles, r.FormValue("role"+value))
			}
		}

		if _, ok := r.MultipartForm.File["img"]; ok {
			file, _, err := r.FormFile("img")

			if err != nil {
				panic(err)
			}

			filename, message, err := service.SaveImage(file, "movie")

			if err != nil {
				panic(err)
			}

			if message != "" {
				tplData := make(html.TemplateData, 0)
				tplData["message"] = message
				html.Render(w, r, "movie/create", tplData)
				return
			}
			movie.ImagePath = filename
		}

		movie.Create()

		http.Redirect(w, r, "/movie/list", http.StatusFound)
		return
	} else {
		panic("Method not allowed")
	}
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	movieId, err := strconv.Atoi(requestVars["id"])

	if err != nil {
		panic(err)
	}

	movie, err := movieModel.GetOne(movieId)

	if err != nil {
		panic(err)
	}

	if r.Method == "GET" {
		actors, err := actorModel.GetAll(0)
		if err != nil {
			panic(err)
		}
		tplData := make(html.TemplateData, 0)
		tplData["movie"] = movie
		tplData["actors"] = actors

		movie.LoadActors()

		html.Render(w, r, "movie/create", tplData)
		return
	} else if r.Method == "POST" {
		r.ParseMultipartForm(MB * 5)

		movie := movieModel.Movie{
			Id:          movieId,
			Name:        strings.Trim(r.FormValue("name"), " "),
			Description: strings.Trim(r.FormValue("description"), " "),
			ImagePath:   movie.ImagePath,
			Year: strings.Trim(r.FormValue("year"), " "),
		}

		ok, message, err := movie.Validate("update")

		if err != nil {
			panic(err)
		}

		if !ok {
			tplData := make(html.TemplateData, 0)
			tplData["message"] = message
			html.Render(w, r, "movie/create", tplData)
			return
		}

		actorIds := r.Form["actors[]"]

		if len(actorIds) > 0 {
			movie.ActorIds = actorIds

			for _, value := range actorIds {
				movie.Roles = append(movie.Roles, r.FormValue("role"+value))
			}
		}

		if _, ok := r.MultipartForm.File["img"]; ok {
			if movie.ImagePath != "" {
				err := os.Remove("public/images/movie/" + movie.ImagePath)
				if err != nil {
					panic(err)
				}
			}

			file, _, err := r.FormFile("img")
			if err != nil {
				panic(err)
			}
			filename, message, err := service.SaveImage(file, "movie")

			if err != nil {
				panic(err)
			}

			if message != "" {
				tplData := make(html.TemplateData, 0)
				tplData["message"] = message
				html.Render(w, r, "movie/create", tplData)
				return
			}

			movie.ImagePath = filename
		}

		movie.Update()

		http.Redirect(w, r, "/movie/list", http.StatusFound)
		return
	} else {
		panic("Method not allowed")
	}
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	actorId, err := strconv.Atoi(requestVars["id"])

	if err != nil {
		panic(err)
	}

	movie, err := movieModel.GetOne(actorId)

	if err != nil {
		panic(err)
	}

	movie.Delete()

	http.Redirect(w, r, "/movie/list", http.StatusFound)
}

func ListMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := movieModel.GetAll(0)
	if err != nil {
		panic(err)
	}
	if r.Method == "GET" {
		tplData := make(html.TemplateData, 0)
		tplData["movies"] = movies
		html.Render(w, r, "movie/list", tplData)
		return
	} else {
		panic("Method not allowed")
	}
}

func MovieDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		requestVars := mux.Vars(r)
		movieId, err := strconv.Atoi(requestVars["id"])

		if err != nil {
			panic(err)
		}

		movie, err := movieModel.GetOne(movieId)

		if err != nil {
			panic(err)
		}

		err = movie.LoadActors()

		if err != nil {
			panic(err)
		}

		movie.GetImdbRating()

		tplData := make(html.TemplateData, 0)
		tplData["movie"] = movie
		html.Render(w, r, "movie/details", tplData)
		return
	} else {
		panic("Method not allowed")
	}
}
