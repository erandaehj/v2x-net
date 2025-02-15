#!/usr/bin/env bash
git switch master
git reset --hard
sudo git clean -fdx
git branch -d local --force
git switch -c local
docker rm -v -f $(docker ps -qa)
sudo systemctl restart docker