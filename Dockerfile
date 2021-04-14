# syntax = docker/dockerfile:1.2
FROM golang:1.16.2-buster

# install dependencies.
RUN apt update && \
    apt-get install -y \
        build-essential \
        ca-certificates \
        curl
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash - && \
    apt-get install -y nodejs

# enable faster module downloading.
ENV GOPROXY https://proxy.golang.org

WORKDIR /starport

# cache dependencies.
COPY ./go.mod . 
COPY ./go.sum . 
RUN go mod download

# build Starport.
COPY . .
# WARNING!
# use `DOCKER_BUILDKIT=1` with `docker build` to enable --mount feature.
RUN --mount=type=cache,target=/root/.cache/go-build go install -v ./...
RUN rm -rf /starport

ENTRYPOINT ["/go/bin/starport"]
