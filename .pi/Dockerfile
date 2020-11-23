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
COPY ipfs.service /etc/systemd/system/ipfs.service
RUN systemctl enable ipfs
