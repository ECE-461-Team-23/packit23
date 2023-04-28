# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /app/

# for web app purposes
EXPOSE 8080

# CMD ["/package_rater_app"]
ENTRYPOINT [ "go", "run", "." ]