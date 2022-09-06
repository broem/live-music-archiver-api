# live-music-archiver-api

This is the API for the [live-music-archiver-extension](https://github.com/broem/live-music-archiver-extension)

## Installation

You must have [Go](https://golang.org/) installed. Version 1.13 or higher is recommended.

You must have [PostgreSQL](https://www.postgresql.org/) installed. Version 12 or higher is recommended.

Clone this repository and run `go build` in the root directory. This will create a binary file called 'live-music-archiver-api'.

## Configuration

Run the 'db.sql' file in your PostgreSQL database. This will create the database and tables. Create a user and grant them access to the database.

Create a file called '.env' in the `server` directory. This file will contain your environment variables. Refer to the `.env.example` file for the required variables.

Update the `config.json` file in the `server` directory. This file contains the configuration for the API.

## Usage

Run the binary file. This will start the API.

## General Information

This API is build using [Echo](https://echo.labstack.com/). 

The database is managed using [go-pg](github.com/go-pg/pg/v10) (this will be changed in the future).

