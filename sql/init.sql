CREATE USER proxy_user WITH password 'qwerty123456';

create database proxy_db
    with owner proxy_user
    encoding 'utf8'
    TABLESPACE = pg_default
;
