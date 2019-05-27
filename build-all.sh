#!/bin/bash


# build our react app, and put files inside of front/build
cd front/
npm install
npm run build

# build docker-stalker backend
cd ../back
export GOARCH=amd64
export GOOS=linux
go build -ldflags="-s -w" .


# build docker image
cd ../
docker build -t 10forward/docker-stalker .

# push docker image
docker push 10forward/docker-stalker






