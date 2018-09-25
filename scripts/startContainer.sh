#!/bin/bash

image=cabal-auth-service-image
container=cabal-auth-service-con

docker stop $container
docker rm $container

docker run -d -p 3000:3000 --name $container $image