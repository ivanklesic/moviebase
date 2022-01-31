package http

import (
	"encoding/json"
	movieModel "moviebase/moviebase/model/entity/movie"
	actorModel "moviebase/moviebase/model/entity/actor"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ApiListMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		requestVars := mux.Vars(r)
		page, err := strconv.Atoi(requestVars["page"])

		if err != nil {
			panic(err)
		}

		movies, err := movieModel.GetAll(page)

		if err != nil {
			panic(err)
		}

		jsonMovies, err := json.Marshal(movies)

		if err != nil {
			panic(err)
		}

		w.Write(jsonMovies)
		return

	} else {
		panic("Method not allowed")
	}
}

func ApiMovieDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		requestVars := mux.Vars(r)
		id, err := strconv.Atoi(requestVars["id"])

		if err != nil {
			panic(err)
		}

		movie, err := movieModel.GetOne(id)

		movie.LoadActors()

		if err != nil {
			panic(err)
		}

		jsonMovie, err := json.Marshal(movie)

		if err != nil {
			panic(err)
		}

		w.Write(jsonMovie)
		return
	} else {
		panic("Method not allowed")
	}

}

func ApiListActors(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		requestVars := mux.Vars(r)
		page, err := strconv.Atoi(requestVars["page"])

		if err != nil {
			panic(err)
		}

		actors, err := actorModel.GetAll(page)

		if err != nil {
			panic(err)
		}

		jsonActors, err := json.Marshal(actors)

		if err != nil {
			panic(err)
		}

		w.Write(jsonActors)
		return

	} else {
		panic("Method not allowed")
	}
}

func ApiActorDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		requestVars := mux.Vars(r)
		id, err := strconv.Atoi(requestVars["id"])

		if err != nil {
			panic(err)
		}

		movie, err := movieModel.GetOne(id)

		movie.LoadActors()

		if err != nil {
			panic(err)
		}

		jsonMovie, err := json.Marshal(movie)

		if err != nil {
			panic(err)
		}

		w.Write(jsonMovie)
		return
	} else {
		panic("Method not allowed")
	}

}
