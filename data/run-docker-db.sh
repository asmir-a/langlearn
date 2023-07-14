docker volume create database
docker run --name database --volume database:/var/lib/postgresql/data --publish 8787:5432 -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=langlearn -d postgres:15-alpine
