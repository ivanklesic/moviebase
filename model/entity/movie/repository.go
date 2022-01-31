package movie

import (
	"moviebase/moviebase/db"
	"moviebase/moviebase/config"
	"strconv"
)

func (movie *Movie) Create() error {
	result, err := db.GetConnection().Exec(
		"insert into movie(title, description_text, image_path, release_year) values(?, ?, ?, ?)",
		movie.Name, movie.Description, movie.ImagePath, movie.Year)

	if err != nil {
		return err
	}

	movieId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	movie.Id = int(movieId)

	if len(movie.ActorIds) > 0 {
		err = movie.persistActors()

		if err != nil {
			return err
		}
	}	
	
	return nil
}

func (movie *Movie) Update() error {
	_, err := db.GetConnection().Exec(
		"update movie set title = ?, description_text = ?, image_path = ?, release_year = ? where id = ?",
		movie.Name, movie.Description, movie.ImagePath, movie.Year, movie.Id)

	if err != nil {
		return err
	}

	err = movie.resetActors()

	if err != nil {
		return err
	}	

	return nil
}

func (movie *Movie) Delete() error {
	_, err := db.GetConnection().Exec(
		"delete from movie where id = ?",
		movie.Id)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetOne(id int) (*Movie, error) {
	var movie Movie
	err := db.GetConnection().QueryRow(
		"select id, title, description_text, image_path, release_year from movie where id = ?",
		id).Scan(&movie.Id, &movie.Name, &movie.Description, &movie.ImagePath, &movie.Year)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func GetAll(page int) (*[]Movie, error) {
	//page == 0 means no pagination

	itemsPerPage := config.GetConfig().Api["itemsPerPage"]

	if page < 0 {
		page = 0
	}

	query := "SELECT id, title, description_text, image_path, release_year FROM movie "

	if page > 0 {
		query += "limit " + strconv.Itoa(itemsPerPage) + " offset " + strconv.Itoa((page - 1) * itemsPerPage)
	}

	rows, err := db.GetConnection().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.Id, &movie.Name, &movie.Description, &movie.ImagePath, &movie.Year); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &movies, nil
}

func Exists(table string, keys []string, values []interface{}) (bool, error) {
	var exists bool
	query := "select COUNT(1) from movie where "
	for index, field := range keys {
		query += field + " = ?"
		if len(keys) > index + 1 {
			query += " and "
		}
	}
	err := db.GetConnection().QueryRow(query, values...).Scan(&exists)

	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (movie *Movie) persistActors() error {
	movieActorQuery := "insert into movie_actor(role_name, movie_id, actor_id) values "
	for i, value := range movie.ActorIds {
		movieActorQuery += "('" + movie.Roles[i] + "', '" + strconv.Itoa(movie.Id) + "', '" + value + "') "
	}

	_, err := db.GetConnection().Exec(movieActorQuery)

	if err != nil {
		return err
	}
	
	return nil
}

func (movie *Movie) resetActors() error {

	_, err := db.GetConnection().Exec("delete from movie_actor where movie_id = ?", strconv.Itoa(int(movie.Id)))

	if err != nil {
		return err
	}
	
	err = movie.persistActors()

	if err != nil {
		return err
	}
	
	return nil
}

func (movie *Movie) LoadActors() error {
	rows, err := db.GetConnection().Query("select actor.id, actor.first_name, actor.last_name, ma.role_name from actor left join movie_actor ma on actor.id = ma.actor_id left join movie on movie.id = ma.movie_id  where movie.id = ?", movie.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var actorId int
		var actorFirstName, actorLastName, role string
		if err := rows.Scan(&actorId, &actorFirstName, &actorLastName, &role); err != nil {
			return err
		}
		movie.ActorIds = append(movie.ActorIds, strconv.Itoa(actorId))
		movie.ActorNames = append(movie.ActorNames, actorFirstName + " " + actorLastName)
		movie.Roles = append(movie.Roles, role)
	}

	if err = rows.Err(); err != nil {
		return  err
	}

	return nil
}