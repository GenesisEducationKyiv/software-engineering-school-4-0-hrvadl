# Mailer microservice

This service is responsible for sending email using 3rd party provider. Currently, I'm using [Resend](https://resend.com/overview) as a mail provider.

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
2. Run `task run` from `./mailer` dir or `task mailer:run` from the root of the repo

## App diagram

<img width="756" alt="image" src="https://github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/assets/93580374/dd9e8930-da1b-41da-a7d5-2c1b4209f701">



## Folder structure

1. `pkg` contains possibly reusable package, not binded to this project. Currently it contains only logger utils
2. `internal`contains packages binded to this project.
   - `cfg` contains config which is read from environment vars.
   - `app` is an abstraction with all services initialization.
   - `transport` contains all transport layer logic: grpc server.
   - `service` contains all domain logic.
3. `cmd` contains entrypoints to the program.
4. `platform` contains specific implementations for sending mails which could possibly change/be changed.
