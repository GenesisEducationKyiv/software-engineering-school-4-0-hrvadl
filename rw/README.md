# Ratewatcher microservice (rw)

This service is responsible for getting latest exchange rate for USD -> UAH and submitting it to the queue.

## Available tasks

You can see all available tasks running following command in the root of the repo:

```sh
task
```

You should get a following output:

```sh
task: [default] task --list-all
task: Available tasks for this project:
* default:               Show available tasks
* generate:              Generate (used for mock generation)
* install:               Install all tools
* lint:                  Run golangci-lint
* run:                   Populate env from .env file and run service
* run-with-env:
* test:                  Run tests
* install:gofumpt:       Install gofumpt
* install:lint:          Install golangci-lint
* install:mock:          Install mockgen
* test:cover:            Run tests & show coverage
* test:race:             Run tests with a race flag
```

## How to run?

If you want to run it as a standalone service you need:

1. Populate env vars needed for it in root `.env` file (../.env)
2. Run `task run` from `./rw` dir or `task rw:run` from the root of the repo

## App diagram

<img width="683" alt="image" src="https://github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/assets/93580374/e458e7d9-c37f-4cd7-bd96-de37e66159d9">

## Folder structure

1. `pkg` contains possibly reusable package, not binded to this project. Currently it contains only logger utils
2. `internal`contains packages binded to this project.

   - `cfg` contains config which is read from environment vars.
   - `app` is an abstraction with all services initialization.
   - `transport` contains all transport layer logic: grpc server & nats publisher.
   - `service` contains all domain logic
   - `platform` contains all platform specific implementations, which can change in the future.
   - `archtest` contains architecture dependency checks.
   - `models` contains application domains models.

3. `cmd` contains entrypoints to the program.
4. `platform` contains specific implementations for querying latest rate exhange, which could change/be changed.
