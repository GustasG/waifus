# Waifus

Frontend interface for [Anime Girls Holding Programming Books](https://github.com/cat-milk/Anime-Girls-Holding-Programming-Books) repository

Made with templ + HTMX

## Prerequisites

Compiling and using this project requires several tools to be installed first

This repository uses git submodules to pull everything from submodules you can use command `git pull --recurse-submodules`

1. Install [Go](https://go.dev)
2. Install [Node.js](https://nodejs.org)

### Setting up

1. Install [templ](https://templ.guide)
2. Install [Task](https://taskfile.dev)
3. (Optionally for live reloading) install [air](https://github.com/cosmtrek/air)

### Building

1. Run `task install` to install necessary dependencies
2. Run `task build` to produce finally binary

### Developing

1. Run `task install` to install necessary dependencies
2. Run `task dev` this will start development server on localhost:3000

### Defining environment variables

1. PORT - Optional variable for which port to use (defaults to 5000)
