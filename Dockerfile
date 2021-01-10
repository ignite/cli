# Build inside lopsided/archlinux for multiplatform support
FROM lopsided/archlinux AS builder
ENV GOPATH=/go
RUN pacman -Syyu --noconfirm go npm make git which && \
	mkdir /go
COPY . /starport
WORKDIR /starport
RUN PATH=$PATH:/go/bin && \
		make

# Copy into a distroless image so that ONLY the starport binary remains
FROM gcr.io/distroless/base
COPY --from=builder /starport/build/starport /

EXPOSE 12345
EXPOSE 8080
EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

CMD ["/starport"]
