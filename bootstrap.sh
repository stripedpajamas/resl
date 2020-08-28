#!/bin/bash
# This is only run when ec2 instance is initially created
# It installs Docker, Node, PM2, and RESL. 

set -e
cd /tmp

# docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker ubuntu

# node
wget https://nodejs.org/dist/v12.18.3/node-v12.18.3-linux-x64.tar.gz
tar -xvf node-v12.8.3-linux-x64.tar.gz
cd node-v12.8.3-linux-x64
sudo cp * /usr/local/ -r

# pm2
npm install -g pm2

# resl
sudo mkdir -p /srv/resl
sudo git clone https://github.com/stripedpajamas/resl.git /srv/resl
pushd /srv/resl
npm install
pm2 start start.js
sudo su -c "env PATH=$PATH:/usr/local/bin pm2 startup -u ubuntu"
popd

