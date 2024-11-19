#!/bin/sh

docker run -p 6379:6379 --name dconf-redis -d redis