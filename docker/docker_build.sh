#!/bin/sh

set -e

BRANCH=$1
PR=$2
TAG=$3

if [ x"" != x$PR ]; then
	echo Skipping docker image build for the pull request
	exit 0
fi

if [ xnext == x$BRANCH ]; then
	IMAGEVERSION=next;
elif [ xmaster == x$BRANCH ]; then
	IMAGEVERSION=latest
else
	echo Don\'t know, how to build for branch $BRANCH
fi

docker build -t akashihi/mdg-rate-loader:$IMAGEVERSION .

if [ x"" != x$TAG ]; then
	docker tag akashihi/mdg-rate-loader$IMAGEVERSION akashihi/mdg-rate-loader:$TAG
fi

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

docker push akashihi/mdg-rate-loader:$IMAGEVERSION

if [ x"" != x$TAG ]; then
	docker push akashihi/mdg-rate-loader:$TAG
fi
