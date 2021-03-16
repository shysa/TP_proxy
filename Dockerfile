#
# go
#
FROM golang:latest

ADD . /d/app
WORKDIR /d/app

RUN go build -o main .

#
# postgresql
#
FROM ubuntu:20.04

ENV PGVER 10
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

USER postgres

ADD ./sql/tables.sql /d/tables.sql
ADD ./sql/init.sql /d/init.sql

# Create a PostgreSQL role named ``dbms_db`` with ``qwerty123456`` as the password and
# then create a database `dbms_db`
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER proxy_user WITH SUPERUSER PASSWORD 'qwerty123456';" &&\
    psql -f /d/init.sql &&\
    psql -f /d/tables.sql -d dbms_db &&\
    /etc/init.d/postgresql stop


RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 512MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "work_mem = 16MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

EXPOSE 8080
EXPOSE 8000

WORKDIR /usr/src/app

COPY . .
COPY --from=0 /d/app/main .

CMD service postgresql start && ./main