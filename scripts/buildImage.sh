#!/bin/bash

echo "Building Docker"
docker rmi -f cabal-auth-service-image
docker build -t cabal-auth-service-image ../ 