# Org Chart Backend

A Go app that gets employee data and sends it to the React frontend.

## Overview

This app:
1. Checks if we already have employee data in our database
2. If not, gets it from the employee API
3. Sends the data to the frontend to show the org chart

## Setup

1. Make sure you have Go installed
2. Clone the repo and go to the backend folder
3. Run `go mod download` to get what you need

## Running it

1. Rename the `.env.example` file into `.env`
1. Set up your database info in the `.env` file
2. Run `go run main.go` or alternatively run `./run.sh` for a fresh start
3. The API will be ready at http://localhost:8080

## API

### GET /employees

Fetches all employees and returns the result in a JSON array format.

## How it works

When you start the app, it:
1. Connects to the database
2. If no employee data is found, gets it from:
   `https://gist.githubusercontent.com/chancock09/6d2a5a4436dcd488b8287f3e3e4fc73d/raw/fa47d64c6d5fc860fabd3033a1a4e3c59336324e/employees.json`
3. Stores the data for next time
4. 