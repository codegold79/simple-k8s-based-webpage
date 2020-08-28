FROM golang:1.15 AS build

WORKDIR /go/src/app
COPY . .

RUN go build -o build/server server.go

FROM photon:3.0

WORKDIR /bin/

COPY --from=build /go/src/app/build .
COPY go-template.html .

ENTRYPOINT [ "server" ]