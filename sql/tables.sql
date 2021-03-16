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
;
GRANT ALL PRIVILEGES ON DATABASE proxy_db TO proxy_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO proxy_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO proxy_user;