# USD -> UAH rate notifier microservice 💸

## Description 💬

The app contains 4 microservices:

- Subscriber (sub) for subscribing users and executing cron job to trigger email once a day at 12:00 UTC (15:00 by Kyiv time)
- Gateway (gw) for mapping HTTP -> GRPC requests and entry point purposes
- Mailer - dumb service for sending emails
- RateWatcher (rw) - service for getting the latest currency exchange rates
- Shared package (pkg) - pkg containing code, which can be used accross modules (logger setup etc.)

As per the task, I need to send a link to only one repository, it was decided to use go workspaces to fit all microservices to one repo. Typically, it should not be the case and it's antipattern. Basically, you can treat each top-level directory as a separate and independent repository/package/module. The `protos` top-level directory is also a go module, containing grpc-generated code.

## Tech stack ⚒️

- [GO](https://go.dev/) as a main programming language
- [`net/http`](https://pkg.go.dev/net/http) for HTTP server
- [GRPC](https://grpc.io/) for inter-service communication
- [`log/slog`](https://go.dev/blog/slog) as a structured logger
- [swaggo](https://github.com/swaggo/swag) for generating swagger documentaion
- [MySQL](https://www.mysql.com/) as a database
- [Golang migrate](https://github.com/golang-migrate/migrate) for DB migrations
- [Gomock](https://github.com/uber-go/mock) for mock generation

## Demo 🎤

https://drive.google.com/file/d/1jrEaqvzlhKB4HB98VEykbiR00G2OEwJS/view?usp=sharing

## How to run? 🏃

1. Copy .env.example contents to .env
2. Get a token for exchange rate API (https://app.exchangerate-api.com/)
3. Populate the `EXCHANGE_API_KEY` variable value with the token you've got
4. Get a token for resend API (https://resend.com/)
5. Verify your domain for sending
6. Populate the `MAILER_API_KEY` variable value with the token you've got
7. Populate the `MAILER_FROM_ADDR` variable with the email you've verified
8. From the root of the repo run `docker compose up -d`

## Local development 🧑🏻‍💻

The repository contains root [taskfile](https://taskfile.dev/) which imports other task files specific to each service. To see all available commands just type:

```sh
task
```

You should get the following output, where each line is the name of the task:

```sh
* default:                      Show available tasks
* lint:
* test:
* gw:generate:                  Generate (used for mock generation)
* gw:install:                   Install all tools
* gw:install:gofumpt:           Install gofumpt
* gw:install:lint:              Install golangci-lint
* gw:install:mock:              Install mockgen
* gw:lint:                      Run golangci-lint
* gw:run:
* gw:test:                      Run tests
* gw:test:cover:                Run tests & show coverage
* gw:test:race:                 Run tests with a race flag
* mailer:generate:              Generate (used for mock generation)
* mailer:install:               Install all tools
* mailer:install:gofumpt:       Install gofumpt
* mailer:install:lint:          Install golangci-lint
* mailer:install:mock:          Install mockgen
* mailer:lint:                  Run golangci-lint
* mailer:run:
* mailer:test:                  Run tests
* mailer:test:cover:            Run tests & show coverage
* mailer:test:race:             Run tests with a race flag
* protos:generate:
* protos:generate:mailer:
* protos:generate:rw:
* protos:generate:sub:
* rw:generate:                  Generate (used for mock generation)
* rw:install:                   Install all tools
* rw:install:gofumpt:           Install gofumpt
* rw:install:lint:              Install golangci-lint
* rw:install:mock:              Install mockgen
* rw:lint:                      Run golangci-lint
* rw:run:
* rw:test:                      Run tests
* rw:test:cover:                Run tests & show coverage
* rw:test:race:                 Run tests with a race flag
* sub:generate:                 Generate (used for mock generation)
* sub:install:                  Install all tools
* sub:install:gofumpt:          Install gofumpt
* sub:install:lint:             Install golangci-lint
* sub:install:mock:             Install mockgen
* sub:lint:                     Run golangci-lint
* sub:run:
* sub:test:                     Run tests
* sub:test:cover:               Run tests & show coverage
* sub:test:race:                Run tests with a race flag
```

They're quite handy to run tests/linters/formatters or to generate grpc-related code and they will install all the deps in case you don't have them.

Before committing anything, you need to also install [left-hook](https://github.com/evilmartians/lefthook). It will run golangci lint before committing to prevent common dummy issues related to formatting/go code. Typically, it should be covered by the CI job, but I need to reduce possible extra runs, considering GH actions has limit on daily runs.

## Compose file 🐋

Compose file has definition of 4 microservices + 1 db service (MySQl) + migrator service (used for running DB migrations on startup). Typically, services won't start till DB is up and migrations has succeeded. So, keep in head, that startup could take some time (1-2 minutes). Data is saved in volume and pods communicate using shared private network. The only pod, which port is mapped to the host port is gateway service.

## Documentation 📄

Services are perfectly documented with Godoc comments. You can host godoc server running following command from the root of the repo:

```sh
task godoc
```

Then you can visit http://localhost:6060/pkg/github.com/hrvadl/converter/?m=all and see documentation for all my packages
<img width="1704" alt="image" src="https://github.com/hrvadl/converter/assets/93580374/e026750c-b399-4801-acc2-76da4d316298">

## App diagram 🏛️

<img width="958" alt="image" src="https://github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/assets/93580374/483a8d8b-d75d-4893-8822-a4e838ea946d">



## CI ⚙️

The application uses GI actions (free tear) as a CI runner. It suits perfectly for small non-commercial projects. CI heavily relies on the [taskfile](https://taskfile.dev/) to do its job. This means each CI step could easily be run locally. CI steps are run in parallel to reduce the time spent waiting for the results.
<img width="1228" alt="image" src="https://github.com/hrvadl/converter/assets/93580374/77b9f5cf-1e9e-485f-a7f8-b29a092811f6">

## Per service documentation 📃

- [Gateway](https://github.com/hrvadl/converter/blob/main/gw/README.md)
- [Mailer](https://github.com/hrvadl/converter/blob/main/mailer/README.md)
- [Ratewatcher](https://github.com/hrvadl/converter/blob/main/rw/README.md)
- [Subscriber](https://github.com/hrvadl/converter/blob/main/sub/README.md)
