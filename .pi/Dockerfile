FROM mydatakeeper/archlinuxarm:aarch64 AS builder
ENV GOPATH=/go
RUN pacman -S --noconfirm go npm make && \
	mkdir /go
COPY . /starport
WORKDIR /starport
RUN PATH=$PATH:/go/bin &&\
		make


FROM faddat/sos

# IPFS systemD service
COPY --from=builder /starport/build/starport /usr/bin/starport
# RUN systemctl enable ipfs
