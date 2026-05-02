FROM golang:1.26 AS builder

WORKDIR /app

RUN apt-get update && \
    apt-get install -y nodejs npm && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN npm install -g @go-task/cli

COPY go.mod go.sum ./
RUN go mod download

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN task build

FROM scratch

EXPOSE 5000

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/assets/css/globals.css ./assets/css/globals.css
COPY --from=builder /app/assets/js/app.min.js ./assets/js/app.min.js
COPY --from=builder /app/assets/favicon.ico ./assets/favicon.ico
COPY --from=builder /app/assets/manifest.json ./assets/manifest.json
COPY --from=builder /app/assets/robots.txt ./assets/robots.txt

ENTRYPOINT ["./main"]
