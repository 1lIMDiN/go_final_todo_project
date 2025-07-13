FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-app .
COPY --from=builder /app/web ./web

COPY --from=builder /app/scheduler.db .

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

VOLUME /data

CMD ["./todo-app"]