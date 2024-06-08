FROM golang:1.22 AS builder

WORKDIR /app

RUN apt update && \
    apt install -y nodejs npm

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN npm install -g @go-task/cli

COPY go.mod go.sum ./
COPY package.json package-lock.json ./
COPY Taskfile.yml .
RUN task install

COPY . .
RUN task build

FROM scratch

EXPOSE 5000

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/assets ./assets

ENTRYPOINT ["./main"]