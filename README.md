# Waifus

Frontend interface for [Anime Girls Holding Programming Books](https://github.com/cat-milk/Anime-Girls-Holding-Programming-Books)

Made with Go, templ, and HTMX. Images are served directly from GitHub via `raw.githubusercontent.com`.

## Prerequisites

1. [Go](https://go.dev)
2. [Node.js](https://nodejs.org)
3. [Task](https://taskfile.dev)
4. [templ](https://templ.guide)
5. [air](https://github.com/cosmtrek/air) — optional, for live reloading

## Setup

### 1. Generate the image manifest

The app reads from `assets/manifest.json` to know which languages and images exist. Generate it by pointing the script at a local clone of the image repository:

```bash
git clone https://github.com/cat-milk/Anime-Girls-Holding-Programming-Books
python scripts/generate_manifest.py Anime-Girls-Holding-Programming-Books -o assets/manifest.json
```

Commit `assets/manifest.json`. Re-run this script whenever the image repository is updated.

### 2. Install dependencies

```bash
task install
```

### 3. Build or run

**Production build:**

```bash
task build
```

**Development server** (live reload on localhost:3000):

```bash
task dev
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `5000` | Port the server listens on |
