# Docs

A microservice skeleton which can help you to build API faster.

## [Micro Service](https://github.com/trisna-ashari/micro)

[![License](https://img.shields.io/badge/License-%20%20Trisnul-red.svg)](https://privy.id)
[![Go](https://img.shields.io/badge/go-1.15-green.svg)](https://golang.org/)
[![Postgres](https://img.shields.io/badge/postgres-13-orange.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/docker-19.03-2885E4.svg)](https://www.docker.com/)

### Table of Contents
- [Getting started](#-getting-started)
- [Preparation](#preparation)
- [Set the env](#set-the-env)
- [Running the service](#running-the-service)
- [Api documentation](#-api-documentation)
- [Available command](#-available-command)
- [Copyright](#-copyright)
-

## 🏃 Getting Started
Hi, in this section, you need to set up the repository in your machine.

👉 Just follow the instruction, and you will in set!

### Preparation
Clone this repo:
```shell script
git clone https://github.com/trisna-ashari/micro.git
```

Download project dependencies:
```shell script
go mod download
go mod tidy
```

Install proto generator:
```shell script
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Set the env
Before you start to run this service, you should set configs with yours.

Create `.env` file. Take a look at: `.env.example`.
```shell script
cp .env.dev.example .env
```

👉  Adjust the credential with installed dependency in your machine.

### Running the service

Run migration:
```shell script
go run main.go db:migrate
```

Run initial seeder:
```shell script
go run main.go db:init
```

Fast run with:
```shell script
go run main.go

# running on default port 6969
```

Run GRPC server:

Open new terminal then run this command

```shell script
go run main.go grpc:start

# running on default port 4949
```

👉 `Hot reload` very useful on development processes

## 📖 Api documentation
This service has builtin API documentation using [swagger](https://github.com/swaggo/swag).
Just run this service and open this link in your browser:

http://localhost:6969/swagger/index.html

To rebuild `api docs`, simply run:
```shell script
swag init
```

## ▶ Available command

This is built-in command in this service:

| Direct Command                 | Build Command   | Description                                                                                            |
|--------------------------------|-----------------|--------------------------------------------------------------------------------------------------------|
| `go run main.go db:migrate`    | `db:migrate`    | Run database migration based on defined struct in `domain/entity` path                                 |
| `go run main.go db:init`       | `db:init`       | Run database seed based on predefined initial factory in `domain/seeds/init_factory.go`                |
| `go run main.go grpc:start`    | `grpc:start`    | Run the GRPC server                                                                                    |

## © Copyright
Trisnul