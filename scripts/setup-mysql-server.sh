#!/bin/bash
docker run --name local-mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql:latest
echo "exec to connect: mysql -h127.0.0.1 -uroot -p123456"