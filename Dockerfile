# Build inside lopsided/archlinux for multiplatform support
# It's an arch linux image with support for both ARM64 and AMD64
FROM lopsided/archlinux


# GOPATH AND GOBIN ON PATH
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin

# INSTALL DEPENDENCIES
RUN pacman -Syyu --noconfirm go npm make git which && \
	mkdir /go

# COPY STARPORT SOURCE CODE INTO CONTAINER
COPY . /starport
WORKDIR /starport

# INSTALL STARPORT
RUN PATH=$PATH:/go/bin && \
		bash scripts/install

# CMD
CMD ["/go/bin/starport"]

# WE NEED BOTH NODE AND GO, DISTROLESS IS NOT THE WAY HERE. REVISIT LATER.
# Copy into a distroless image so that ONLY the starport binary remains
# FROM gcr.io/distroless/base
# COPY --from=builder /starport/build/starport /

# EXPOSE 12345
# EXPOSE 8080
# EXPOSE 1317
# EXPOSE 26656
# EXPOSE 26657

# CMD ["/starport"]
