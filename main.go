package main

import (
	"net/http"

	myhttp "moviebase/moviebase/http/controller"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	//default
	r.HandleFunc("/", myhttp.Index)
	//api
	r.HandleFunc("/api/actor/details/{id}", myhttp.ApiActorDetails)
	r.HandleFunc("/api/actor/list/{page}", myhttp.ApiListActors)
	r.HandleFunc("/api/movie/details/{id}", myhttp.ApiMovieDetails)
	r.HandleFunc("/api/movie/list/{page}", myhttp.ApiListMovies)
	//actors
	r.HandleFunc("/actor/create", myhttp.CreateActor)
	r.HandleFunc("/actor/delete/{id}", myhttp.DeleteActor)
	r.HandleFunc("/actor/update/{id}", myhttp.UpdateActor)
	r.HandleFunc("/actor/details/{id}", myhttp.ActorDetails)
	r.HandleFunc("/actor/list", myhttp.ListActors)
	//movies
	r.HandleFunc("/movie/create", myhttp.CreateMovie)
	r.HandleFunc("/movie/delete/{id}", myhttp.DeleteMovie)
	r.HandleFunc("/movie/update/{id}", myhttp.UpdateMovie)
	r.HandleFunc("/movie/details/{id}", myhttp.MovieDetails)
	r.HandleFunc("/movie/list", myhttp.ListMovies)
	//static file server
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	http.ListenAndServe("localhost:8080", r)
}
