FROM golang:1.16-alpine

RUN apk add --no-cache git && apk add --no-cache bash

WORKDIR /app/GurkhaFabricAPI

COPY go.mod . 
COPY go.sum . 

RUN go mod download

COPY . .

RUN go build -o ./out/GurkhaFabricAPI .

EXPOSE 8080

ENTRYPOINT ["./out/GurkhaFabricAPI"]