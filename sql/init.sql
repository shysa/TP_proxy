CREATE USER proxy_user WITH password 'qwerty123456';

create database proxy_db
    with owner proxy_user
    encoding 'utf8'
    TABLESPACE = pg_default
;

create table if not exists request
(
    id             serial primary key,
    method         text not null,
    scheme         text not null,
    path           text not null,
    proto          text not null,
    host           text not null,
    url            text not null,
    request_text   text not null
);

create table if not exists response
(
    id            serial primary key,
    response_text text not null,
    request_id    integer references request
);

GRANT ALL PRIVILEGES ON database proxy_db TO proxy_user;
