#!/bin/bash
cp ~/.ssh/id_rsa.pub docker/
docker-compose build
docker-compose up -d
