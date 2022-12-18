FROM golang:alpine3.14 as compiler

RUN apk add git

WORKDIR /app/build

COPY . .

RUN go build

FROM alpine:3.14

RUN wget -L --fail https://github.com/docker/compose/releases/download/1.24.0/run.sh -o /usr/local/bin/docker-compose

RUN chmod +x /usr/local/bin/docker-compose

WORKDIR /app/prod

COPY --from=compiler /app/build/whaleman .

ENTRYPOINT ["./whaleman"]