
FROM archlinux AS builder
ENV GOPATH=/go
RUN pacman -S --noconfirm go npm && \
	mkdir /go
COPY . /starport
WORKDIR /starport
RUN make



FROM faddat/sos

RUN pacman -Syyu --noconfirm neovim \
				colordiff \
				tree \
				unzip \
				unrar \
				htop \
				nmap \
				iftop \
				iotop \
				strace \
				lsof \
				git \
				jshon \
				nano \
				bc \
				wget \
				go-ipfs \
				jq \
				go


# Add development environment dependencies
RUN echo "zerotier-cli join 35c192ce9bf02045" >> /usr/local/bin/firstboot.sh && \
		echo "ipfs init" /usr/local/bin/firstboot.sh

# IPFS systemD service
COPY --from=builder /starport/build/starport /usr/bin/starport
# RUN systemctl enable ipfs
