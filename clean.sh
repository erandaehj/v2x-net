#!/usr/bin/env bash
git switch master
git reset --hard
sudo git clean -fdx
docker rm -v -f $(docker ps -qa)
docker network rm fabricnet_test
sudo systemctl restart docker