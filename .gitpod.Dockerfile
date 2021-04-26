FROM gitpod/workspace-full

RUN brew install gh

RUN npm install -g npm@7.10.0
