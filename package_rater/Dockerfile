# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

# copy individually so that if there is cache of any file/folder, line can be skipped
# COPY run ./
# COPY url.txt ./
# COPY cmd ./cmd
# COPY github_apis ./github_apis
# COPY internal ./internal
# COPY tests ./tests

COPY . /app/

# unsure how to build this app
# RUN go build -o /package_rater_app

# for web app purposes
# EXPOSE 8080

# CMD ["/package_rater_app"]