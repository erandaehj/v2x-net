#!/usr/bin/env bash
git switch master
git reset --hard
git branch -d local --force
git switch -c local
docker rm -v -f $(docker ps -qa)
sudo systemctl restart docker