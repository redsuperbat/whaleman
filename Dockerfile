FROM golang:alpine3.14 as compiler

RUN apk add git

WORKDIR /app/build

COPY . .

RUN go build

FROM alpine:3.14

RUN wget https://github.com/docker/compose/releases/download/1.24.0/run.sh -O /usr/local/bin/docker-compose

RUN chmod +x /usr/local/bin/docker-compose

RUN apt update

RUN apt install docker.io

WORKDIR /app/prod

COPY --from=compiler /app/build/whaleman .

ENTRYPOINT ["./whaleman"]