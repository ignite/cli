FROM faddat/sos-full

# Make SPN user, spn hostname so we get spn.local via mdns
RUN useradd --create-home starport && \
        echo "starport" > /etc/hostname && \
        echo "starport:starport" | chpasswd && \
        usermod --append --groups wheel starport && \
        passwd -e root


# later, update to tendermint/starport
COPY --from=tendermintdevelopment/starport /go/bin/starport /usr/bin/starport
