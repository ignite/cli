FROM gitpod/workspace-full

RUN wget https://golang.org/dl/go1.16rc1.linux-amd64.tar.gz && \
    rm -rf $HOME/go && \
    tar -C $HOME -xzf go1.16rc1.linux-amd64.tar.gz && \
    rm go1.16rc1.linux-amd64.tar.gz

ENV GONAME go1.15

RUN brew install gh protobuf
