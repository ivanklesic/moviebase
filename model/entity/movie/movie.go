package movie

import (
	"moviebase/moviebase/config"
	"moviebase/moviebase/service"
	"strconv"
)

type Movie struct {
	Id          int `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImagePath   string `json:"-"`
	Year        string `json:"year"`
	ActorIds    []string `json:"actorIds"`
	ActorNames  []string `json:"-"`
	Roles       []string `json:"-"`
	Rating float64 `json:"-"`
}

type imdbMovieJson struct {
	Results []struct {
		Id string `json:"id"`
	} `json:"results"`
}

type imdbRatingJson struct {
	Rating string `json:"imDb"`
}

func (movie Movie) Validate(method string) (bool, string, error) {
	actorExists, err := Exists("actor", []string{"title"}, []interface{}{movie.Name})
	if err != nil {
		return false, "", err
	}
	if len(movie.Name) > 255 || len(movie.Description) > 10000 {
		return false, "Name or description is too long", nil
	}

	if method == "create" {
		if actorExists {
			return false, "Movie with the same name already exists", nil
		}
	}

	return true, "", nil
}

func (movie *Movie) GetImdbRating() error {
	imdbConfig := config.GetConfig().Imdb
	imdbApiKey := imdbConfig["apiKey"]

	searchString := movie.Name

	if movie.Year != "" {
		searchString += " " + movie.Year
	}

	var imdbMovie imdbMovieJson

	getIdByNameUrl := "https://imdb-api.com/API/SearchMovie/" + imdbApiKey + "/" + searchString

	err := service.GetJson(getIdByNameUrl, &imdbMovie)

	if err != nil {
		return err
	}

	var imdbRating imdbRatingJson

	getRatingByIdUrl := "https://imdb-api.com/API/Ratings/" + imdbApiKey + "/" + imdbMovie.Results[0].Id

	err = service.GetJson(getRatingByIdUrl, &imdbRating)

	if err != nil {
		return err
	}

	rating, err := strconv.ParseFloat(imdbRating.Rating, 64)

	if err != nil {
		return err
	}

	movie.Rating = rating

	return nil
}
