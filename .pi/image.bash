#!/bin/bash

# This install script runs inside the chroot of your image builder.
# After it runs, a second shell session installs nodejs with nvm

# Connect to the internet by configuring DNS
mv /etc/resolv.conf /etc/resolv.conf.bk
echo 'nameserver 8.8.8.8' > /etc/resolv.conf
echo 'nameserver 1.1.1.1' >> /etc/resolv.conf

# Update packages and install Zerotier
apt update
apt install -y curl wget jq apt-transport-https gnupg gnupg-agent software-properties-common
wget http://download.zerotier.com/debian/buster/pool/main/z/zerotier-one/zerotier-one_1.4.6_arm64.deb
dpkg -i zerotier-one_1.4.6_arm64.deb
apt install -y

# Set Up Starport
chmod +x /usr/bin/starport
mkdir -p /starport/config
systemctl enable starportd

# Set up Nodejs without NVM
wget https://nodejs.org/dist/v12.18.4/node-v12.18.4-linux-arm64.tar.xz
sudo mkdir -p /usr/local/lib/nodejs
sudo tar -xJvf node-v12.18.4-linux-arm64.tar.xz -C /usr/local/lib/nodejs
echo "export PATH=/usr/local/lib/nodejs/node-v12.18.4-linux-arm64/bin:$PATH" >> /root/.profile
echo "export PATH=/usr/local/lib/nodejs/node-v12.18.4-linux-arm64/bin:$PATH" >> /home/ubuntu/.profile
. /root/.profile
node -v
npm version
npx -v
