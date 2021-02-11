FROM gitpod/workspace-full

RUN brew install gh protobuf git-lfs

RUN git lfs install
