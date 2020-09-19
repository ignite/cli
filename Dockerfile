FROM archlinux:latest

# COPY . .
RUN pacman -Syyu --noconfirm
RUN pacman -S --noconfirm go nodejs yarn git npm base-devel make
RUN git clone https://github.com/tendermint/starport
# RUN make
RUN npm i -g @tendermint/starport
