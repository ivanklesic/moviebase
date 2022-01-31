drop database moviebase;

create database moviebase;
alter database moviebase character set utf8mb4 collate utf8mb4_unicode_ci;

use moviebase;

create table actor (
    id int auto_increment primary key not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    image_path varchar(255) default ''
);

create table movie(
    id int auto_increment primary key not null,
    title varchar(255),
    release_year varchar(4),
    description_text text,
    image_path varchar(255) default ''
);

create table movie_actor(
    id int auto_increment primary key not null,
    role_name varchar(255) default null,
    actor_id int not null,
    movie_id int not null,
    FOREIGN KEY (actor_id) REFERENCES actor(id)
        on delete cascade,
    FOREIGN KEY (movie_id) REFERENCES movie(id)
        on delete cascade
);