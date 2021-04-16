# syntax = docker/dockerfile:1.2
# WARNING! Use `DOCKER_BUILDKIT=1` with `docker build` to enable --mount feature.

## prep the base image.
#
FROM golang:1.16.2-buster as base

RUN apt update && \
    apt-get install -y \
        build-essential \
        ca-certificates \
        curl

RUN curl -sL https://deb.nodesource.com/setup_14.x | bash - && \
    apt-get install -y nodejs

# enable faster module downloading.
ENV GOPROXY https://proxy.golang.org

## builder stage.
#
FROM base as builder

WORKDIR /starport

# cache dependencies.
COPY ./go.mod . 
COPY ./go.sum . 
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build go install -v ./...

## prep the final image.
#
FROM base

COPY --from=builder /go/bin/starport /usr/bin

# see docs for exposed ports:
#   https://docs.starport.network/4%20Configure%20a%20Blockchain/2%20config.yml%20Reference.html#host
EXPOSE 26657
EXPOSE 26656
EXPOSE 6060 
EXPOSE 9090 
EXPOSE 1317 
EXPOSE 8080
EXPOSE 12345 

ENTRYPOINT ["starport"]
