#!/bin/bash -e
apt-get -y update
apt-get -y install make curl bison
apt-get -y install git mercurial 

# install go1.4 using gvm
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /home/vagrant/.gvm/scripts/gvm
gvm install go1.4
gvm use go1.4

cat >> /home/vagrant/.bashrc <<EOF
gvm use go1.4
EOF

# install gpm https://github.com/pote/gpm
git clone https://github.com/pote/gpm.git && cd gpm
git checkout v1.3.1
./configure
make install
cd ..

# install gvp
git clone https://github.com/pote/gvp.git && cd gvp
git checkout v0.1.0
./configure
make install
cd ..


