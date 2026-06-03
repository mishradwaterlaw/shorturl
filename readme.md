# URL Shortener API

A note about this project: this was built specifically for learning Go, not as a portfolio project.

If you're reading this, the goal was to work through Go concepts by building something real instead of only following abstract tutorials.

## What it does

This project is a simple URL shortener HTTP API.

- Send a `POST` request with a long URL.
- Get a short code back.
- Send a `GET` request with that short code.
- Get redirected to the original URL.

## What I learned building this

Working on this project helped reinforce a number of core Go concepts:

- Go modules and package structure
- Structs and zero values
- Multiple return values and error handling
- Interfaces satisfied implicitly
- Pointer receivers
- Goroutines and channels
- `sync.RWMutex` for concurrent map access
- Table-driven tests
- Go's standard library HTTP server

## Note on comments

The inline code comments were AI-assisted to help explain and revise concepts during the learning process.

The code itself was written by hand through a guided learning session.

## How to run it

```bash
go run main.go
```

Then send a `POST` request to `http://localhost:8080/shorten` with JSON like:

```json
{"url": "your-url-here"}
```