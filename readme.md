# Go Tic-Tac-Toe

A web application featuring multiplayer online tic-tac-toe battle.

## Installation
Have the following installed for a working development environment:
* [GoLang](https://github.com/golang/go) - Written in
* [PostgresSQL](https://www.postgresql.org/) - Relational Database
* [Goose](https://github.com/pressly/goose) - Database Schema Migration Tool

## Configuration
In the root directory, create a `config.toml` and fill it up the following:

`[database]`
`type = "postgres"`
`name = "DATABASE_NAME"`
`user = "DATABASE_USER"`
`host = "localhost"`
`port = "PORT_NO"`
`sslmode = "disable"`
`[server]`
`port = "SERVER_PORT_NO"`


Replace quotes in caps with the relevant information.

## Running the test suite
Enter `main.go` and run `go test`. 


