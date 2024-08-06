FROM golang:1.22 AS builder

WORKDIR /app

RUN apt update && \
    apt install -y nodejs npm && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN npm install -g @go-task/cli

COPY go.mod go.sum package.json package-lock.json Taskfile.yml ./
RUN task install

COPY . .
RUN task build

FROM scratch

EXPOSE 5000

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/assets ./assets

ENTRYPOINT ["./main"]