#!/usr/bin/env bash
git checkout develop
git branch -d local
git checkout local
docker rm -v -f $(docker ps -qa)
sudo systemctl restart docker