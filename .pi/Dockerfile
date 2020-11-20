FROM faddat/spos

RUN pacman -Syyu --noconfirm base-devel \
				nvim \
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

# INSTALL YAY: WE BECOME THE BUILDUSER HERE
USER builduser
RUN cd ~/ && \
		git clone https://aur.archlinux.org/yay.git yay && \
		cd yay && \
		makepkg -si --noconfirm --clean --rmdeps
USER root

# Add development environment dependencies
RUN echo "zerotier-cli join 35c192ce9bf02045" >> /usr/local/bin/firstboot.sh && \
		echo "ipfs init" /usr/local/bin/firstboot.sh


# IPFS systemD service
COPY contrib/ipfs.service /etc/systemd/system/ipfs.service
RUN systemctl enable ipfs





