#!/usr/bin/env bash
git switch master
git reset --hard
sudo git clean -fd
git branch -d local --force
git switch -c local
docker rm -v -f $(docker ps -qa)
sudo systemctl restart docker