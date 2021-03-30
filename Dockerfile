FROM golang:1.16.2-alpine3.13

# INSTALL DEPENDENCIES
RUN apk add --no-cache npm make git bash which curl protoc

# INSTALL PROTOBUF LIBRARY
ENV PROTOC_ZIP=protoc-3.13.0-linux-x86_64.zip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/${PROTOC_ZIP}
RUN unzip -o ${PROTOC_ZIP} -d /proto
RUN cp -R /proto/include/* /usr/include/

# COPY STARPORT SOURCE CODE INTO CONTAINER
COPY . /starport
WORKDIR /starport

# INSTALL STARPORT
RUN make install

# CMD
ENTRYPOINT ["/go/bin/starport"]

# EXPOSE PORTS
EXPOSE 12345
EXPOSE 8080
EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

