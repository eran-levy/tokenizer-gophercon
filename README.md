### Tokenizer gophercon sample service
Setup pre-requisite - 
go 1.15+
docker (desktop + kuberneres / other k8s distribution)
skaffold
helm

setup and start (scripts/):
1. mysql
2. redis
3. prometheus
4. jaeger all-in-one
5. protoc compiler and go plugins
6. execute sql to create the database and table

Jaeger UI: http://localhost:16686/search

Prometheus UI: http://localhost:9090/

Redis cli:

`redis-cli
127.0.0.1:6379> ping
PONG
127.0.0.1:6379> select 0
OK
127.0.0.1:6379> keys *
(empty array)`

Service environment variables override (DSN env var: SERVICE_DB_DSN=root:123456@tcp(127.0.0.1:3306)/tokenizer):

`SERVICE_LOG_LEVEL=debug;SERVICE_DB_PASSWD=123456;`


