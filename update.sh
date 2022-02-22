#!/bin/bash

cd /home/pi

rm -rf /home/pi/CUBEIOT
git clone https://github.com/Kae-Tempest/CUBEIOT

cd CUBEIOT/frontend

go build .

cd /home/pi/CUBEIOT/backend

sudo chmod 700 ./install-venv.sh
./install-venv.sh

systemctl restart web
systemctl restart api