FROM golang:alpine3.14 as compiler

RUN apk add git

WORKDIR /app/build

COPY . .

RUN go build

FROM linuxserver/docker-compose:1.29.2-alpine

WORKDIR /app/prod

COPY --from=compiler /app/build/whaleman .

ENTRYPOINT ["./whaleman"]