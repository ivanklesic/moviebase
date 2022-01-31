package actor

import (
	"moviebase/moviebase/config"
	"moviebase/moviebase/db"
	"strconv"
)

func (actor *Actor) Create() error {
	result, err := db.GetConnection().Exec(
		"insert into actor(first_name, last_name, image_path) values(?, ?, ?)",
		actor.FirstName, actor.LastName, actor.ImagePath)

	if err != nil {
		return err
	}

	actorId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	actor.Id = int(actorId)

	if len(actor.MovieIds) > 0 {
		err = actor.persistMovies()

		if err != nil {
			return err
		}
	}	
	
	return nil
}

func (actor *Actor) Update() error {
	_, err := db.GetConnection().Exec(
		"update actor set first_name = ?, last_name = ?, image_path = ? where id = ?",
		actor.FirstName, actor.LastName, actor.ImagePath, actor.Id)

	if err != nil {
		return err
	}

	err = actor.resetMovies()

	if err != nil {
		return err
	}	

	return nil
}

func (actor *Actor) Delete() error {
	_, err := db.GetConnection().Exec(
		"delete from actor where id = ?",
		actor.Id)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetOne(id int) (*Actor, error) {
	var actor Actor
	err := db.GetConnection().QueryRow(
		"select id, first_name, last_name, image_path from actor where id = ?",
		id).Scan(&actor.Id, &actor.FirstName, &actor.LastName, &actor.ImagePath)
	if err != nil {
		return nil, err
	}
	return &actor, nil
}

func GetAll(page int) (*[]Actor, error) {
	//page == 0 means no pagination

	itemsPerPage := config.GetConfig().Api["itemsPerPage"]

	if page < 0 {
		page = 0
	}

	query := "select * from actor "

	if page > 0 {
		query += "limit " + strconv.Itoa(itemsPerPage) + " offset " + strconv.Itoa((page - 1) * itemsPerPage)
	}

	rows, err := db.GetConnection().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []Actor

	for rows.Next() {
		var actor Actor
		if err := rows.Scan(&actor.Id, &actor.FirstName, &actor.LastName, &actor.ImagePath); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &actors, nil
}

func Exists(table string, keys []string, values []interface{}) (bool, error) {
	var exists bool
	query := "select COUNT(1) from actor where "
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

func (actor *Actor) persistMovies() error {
	movieActorQuery := "insert into movie_actor(role_name, actor_id, movie_id) values "
	for i, value := range actor.MovieIds {
		movieActorQuery += "('" + actor.Roles[i] + "', '" + strconv.Itoa(actor.Id) + "', '" + value + "') "
	}

	_, err := db.GetConnection().Exec(movieActorQuery)

	if err != nil {
		return err
	}
	
	return nil
}

func (actor *Actor) resetMovies() error {

	_, err := db.GetConnection().Exec("delete from movie_actor where actor_id = ?", strconv.Itoa(int(actor.Id)))

	if err != nil {
		return err
	}
	
	err = actor.persistMovies()

	if err != nil {
		return err
	}
	
	return nil
}

func (actor *Actor) LoadMovies() error {
	rows, err := db.GetConnection().Query("select movie.id, movie.title, ma.role_name from movie left join movie_actor ma on movie.id = ma.movie_id left join actor on actor.id = ma.actor_id  where actor.id = ?", actor.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var movieId int
		var movieName, role string
		if err := rows.Scan(&movieId, &movieName, &role); err != nil {
			return err
		}
		actor.MovieIds = append(actor.MovieIds, strconv.Itoa(movieId))
		actor.MovieNames = append(actor.MovieNames, movieName)
		actor.Roles = append(actor.Roles, role)
	}

	if err = rows.Err(); err != nil {
		return  err
	}

	return nil
}

func ApiGetAll(){

}

func ApiDetails(){
	
}

