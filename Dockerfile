FROM archlinux:latest

COPY .

RUN pacman -Syyu --no-confirm
RUN pacman -S go node yarn git npm base-devel make
RUN make
