#!/bin/bash
docker run --name local-redis -p 6379:6379 -d redis
echo "exec to connect: redis-cli"