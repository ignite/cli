FROM faddat/sos-full

# Make SPN user, spn hostname so we get spn.local via mdns
RUN useradd --create-home starport && \
        echo "starport" > /etc/hostname && \
        echo "starport:starport" | chpasswd && \
        usermod --append --groups wheel starport && \
        passwd -e root

USER starport

# Setup $PATH & GOPATH for more seamlessness
RUN mkdir /home/starport/go && \
        echo "PATH=$PATH:/home/starport/go/bin" >> /home/starport/.bash_profile && \
        echo "GOPATH=/home/starport/go"  >> /home/starport/.bash_profile

USER root


# later, update to tendermint/starport
COPY --from=tendermintdevelopment/starport /go/bin/starport /usr/bin/starport
