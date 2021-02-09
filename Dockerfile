# Alpine provides the same multiplatform support as arch.
FROM alpine

# GOPATH AND GOBIN ON PATH
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin

# INSTALL DEPENDENCIES
RUN apk add --no-cache go npm make git which && \
	mkdir /go

# COPY STARPORT SOURCE CODE INTO CONTAINER
COPY . /starport
WORKDIR /starport

# INSTALL STARPORT
RUN PATH=$PATH:/go/bin && \
		bash scripts/install

# CMD
CMD ["/go/bin/starport"]

# EXPOSE PORTS
EXPOSE 12345
EXPOSE 8080
EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

