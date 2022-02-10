package main

import (
	"net/http"
	"os"

	myhttp "moviebase/moviebase/http/controller"

	"moviebase/moviebase/middleware"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	f, err := os.OpenFile("logs/log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		panic(err)
	}

	//defer f.Close()

	logger := middleware.LoggerInit(f)

	loggingMiddleware := middleware.LoggingMiddleware(logger)

	//default
	router.HandleFunc("/", myhttp.Index)
	//api
	router.HandleFunc("/api/actor/details/{id}", myhttp.ApiActorDetails)
	router.HandleFunc("/api/actor/list/{page}", myhttp.ApiListActors)
	router.HandleFunc("/api/movie/details/{id}", myhttp.ApiMovieDetails)
	router.HandleFunc("/api/movie/list/{page}", myhttp.ApiListMovies)
	//actors
	router.HandleFunc("/actor/create", myhttp.CreateActor)
	router.HandleFunc("/actor/delete/{id}", myhttp.DeleteActor)
	router.HandleFunc("/actor/update/{id}", myhttp.UpdateActor)
	router.HandleFunc("/actor/details/{id}", myhttp.ActorDetails)
	router.HandleFunc("/actor/list", myhttp.ListActors)
	//movies
	router.HandleFunc("/movie/create", myhttp.CreateMovie)
	router.HandleFunc("/movie/delete/{id}", myhttp.DeleteMovie)
	router.HandleFunc("/movie/update/{id}", myhttp.UpdateMovie)
	router.HandleFunc("/movie/details/{id}", myhttp.MovieDetails)
	router.HandleFunc("/movie/list", myhttp.ListMovies)
	//static file server
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	loggingRouter := loggingMiddleware(router)

	if err := http.ListenAndServe("localhost:8080", loggingRouter); err != nil {
		logger.Log("status", "fatal", "err", err)
		os.Exit(1)
	}
}
