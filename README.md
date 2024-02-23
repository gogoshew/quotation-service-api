# Quotation Service API

## Task
_Implement a service for exchange rate quotes.
The service provides an asynchronous interface, i.e. user first
perform a request to update quotes, and then, after some time,
request for quotation. In this case, the update itself
quotes happen in the background._

## Description
Our service represents API for check and update quotations
There are only 3 allowed currencies, base currency is USD it means that all quotation values calculating as, for example, USD/EUR and so on:
- EUR
- MXN
- GEL

Quotes are updated in the background using the Cron package once at a certain point in time specified in the docker-compose environment variables


## General packages

Routing:
- gorilla/mux

DB:
- gorm (pgsql driver)

CRON:
- robfig/cron/v3


## Documentation
All API methods described in `swagger.yaml`

Also, you could use prepared requests in `http` folder

## How to run app

There is `Makefile` in current project:
- pull app images `make build-app`
- run app  `make compose-up` from project root directory
- stop app `make compose-down` from project root directory
