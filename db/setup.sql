create table actor (
    id int auto_increment primary key not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null
);

create table movie(
    id int auto_increment primary key not null,
    name varchar(255),
    release_date date not null,
    won_oscar tinyint not null,
    director_id int not null,
    total_ratings_sum int default 0,
    total_ratings_count int default 0,
    FOREIGN KEY (director_id) REFERENCES director(id)
        on delete cascade
);

create table movie_actor(
    id int auto_increment primary key not null,
    actor_id int,
    movie_id int,
    role varchar(255),
    FOREIGN KEY (actor_id) REFERENCES actor(id)
    on delete cascade,
    FOREIGN KEY (movie_id) REFERENCES movie(id)
        on delete cascade

);

create table reviewer(
    id int auto_increment primary key not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null
);

create table movie_rating(
    id int auto_increment primary key not null,
    reviewer_id int,
    movie_id int,
    rating float,
    FOREIGN KEY (reviewer_id) REFERENCES reviewer(id)
        on delete cascade,
    FOREIGN KEY (movie_id) REFERENCES movie(id)
        on delete cascade
);
