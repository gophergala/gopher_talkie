#!/bin/bash -e
sudo apt-get -y update
sudo apt-get -y install make curl bison
sudo apt-get -y install git mercurial 
sudo apt-get -y install pkg-config
sudo apt-get -y install gnupg
sudo apt-get -y install portaudio19-dev
sudo apt-get -y install sqlite3

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
sudo make install
cd ..

# install gvp
git clone https://github.com/pote/gvp.git && cd gvp
git checkout v0.1.0
./configure
sudo make install
cd ..


