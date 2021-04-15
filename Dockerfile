FROM golang:1.16.2-buster

RUN apt update && \
    apt-get install -y \
        build-essential \
        ca-certificates \
        unzip \
        curl

# INSTALL NODE
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash - && \
    apt-get install -y nodejs

# COPY STARPORT SOURCE CODE INTO CONTAINER
COPY ./docs /starport/docs
COPY ./starport /starport/starport
COPY ./go.mod /starport/go.mod
COPY ./go.sum /starport/go.sum
WORKDIR /starport

# INSTALL STARPORT
RUN go install -mod=readonly ./...

# ENTRYPOINT
ENTRYPOINT ["/go/bin/starport"]

# EXPOSE PORTS
EXPOSE 12345
EXPOSE 8080
EXPOSE 1317
EXPOSE 26656
EXPOSE 26657
