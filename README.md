# stud-api

A lean student service written in Go.

## Overview

`stud-api` is a minimal REST API for managing students.
Built with a standard-library-first mindset.

No frameworks.
No abstractions for the sake of abstractions.

## Features

- CRUD operations for students
- JSON over HTTP
- Explicit request handling
- Proper HTTP status codes
- Simple and predictable routing


## Running Locally

go mod tidy  
go run ./cmd/stud-api

Server runs on:

http://localhost:5000

## Philosophy

- Simple > clever
- Explicit > magical
- Code should explain itself

## Status

Work in progress.
Focused on correctness and readability.
