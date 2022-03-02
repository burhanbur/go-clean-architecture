
# Description
This project is an example of implementation RESTful API blog service. By implementing a Clean Architecture in Go (Golang), it is expected to be a learning tool for the future.

This project has 4 Domain layer :

- Models Layer
  Encapsulate enterprise wide business rules. An entity in Go is a set of data structures and functions.
- Repository Layer
  This layer is a set of adapters that convert data from the format most convenient for the use cases and entities, to the format most convenient for some external agency such as the Database or the Web
- Service Layer
  This layer contains application specific business rules. It encapsulates and implements all of the use cases of the system.
- Controller Layer
  This layer is generally composed of frameworks and tools such as the Database, the Web Framework, etc.

> Each request must be validated by a JWT (JSON Web Token) which is stored in redis once the login endpoint is reached. Expired token duration can be set in .env file

## Installation

- clone apps and install the dependencies

```sh
git https://github.com/burhanbur/laravel-design-pattern.git
cd src
composer update
```

- Configure database and redis in folder config
- Open the program via the terminal.
- Run the program by typing the command

```sh
go run main.go
```
- Rename file .env.example to .env

## What's Inside

* Gin as Go framework
* GORM as ORM library for Go
* MySQL as database driver
* Redis as in-memory cache database
* JWT as token based authentication
* Personalized system logger