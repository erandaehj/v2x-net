#!/usr/bin/env bash
git switch master
git reset --hard
sudo git clean -fdx
docker rm -v -f $(docker ps -qa)
sudo systemctl restart docker