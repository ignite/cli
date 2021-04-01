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

# INSTALL PROTOBUF LIBRARY
RUN curl -sL https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip -o protoc.zip && \
    unzip protoc.zip -d /usr/local && \
    rm protoc.zip

# COPY STARPORT SOURCE CODE INTO CONTAINER
COPY . /starport
WORKDIR /starport

# INSTALL STARPORT
RUN make install

# ENTRYPOINT
ENTRYPOINT ["/go/bin/starport"]

# EXPOSE PORTS
EXPOSE 12345
EXPOSE 8080
EXPOSE 1317
EXPOSE 26656
EXPOSE 26657