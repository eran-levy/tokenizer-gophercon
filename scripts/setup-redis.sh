#!/bin/bash
docker run --name local-redis -d redis
echo "exec to connect: redis-cli"