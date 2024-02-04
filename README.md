# Go Todo

A simple mini-project to learn a little about the following stack:

- Go lang
- Gin
- Gorm
- HTMX

## How to run

The simplest way to run is to use `docker compose`

`docker compose up`

And access `http://localhost:8080/`

## Compiling the backend

Need to have the [just](https://github.com/casey/just) tool installed

`just build`

## Configuration

The configuration of the backend is done on the `.env` file. Just refer to the `.env.example` file.

### SQLite storage

```
DB_ENGINE=SQLITE
DB_FILE=todo.db
```

### PostgreSQL storage

```
DB_ENGINE=POSTGRES
DB_DSN=host=postgres user=usr password=pw sslmode=disable
```
