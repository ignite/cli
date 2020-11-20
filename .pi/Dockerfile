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


# Test 1



