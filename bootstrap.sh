#!/bin/bash
# This is only run when ec2 instance is initially created
# It installs Docker, Node, PM2, and RESL. 

set -e
cd /tmp

# docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
usermod -aG docker ubuntu

# node
wget https://nodejs.org/dist/v12.18.3/node-v12.18.3-linux-x64.tar.gz
tar -xvf node-v12.18.3-linux-x64.tar.gz
cd node-v12.18.3-linux-x64
cp * /usr/local/ -r

# pm2
npm install -g pm2

# resl
mkdir -p /srv/resl
git clone https://github.com/stripedpajamas/resl.git /srv/resl
pushd /srv/resl
npm install
pm2 start start.js
env PATH=$PATH:/usr/local/bin /usr/local/lib/node_modules/pm2/bin/pm2 startup systemd -u ubuntu --hp /home/ubuntu
popd

