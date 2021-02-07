# tokenizer-gophercon

Setup pre-requisite - 
docker
microk8s

docker start jaeger

UI: http://localhost:16686/search
apps/prometheus-2.10.0.darwin-amd64 ./prometheus --config.file=prometheus.yml --log.level=error
UI: http://localhost:9090/
histogram_quantile(0.95, sum(rate(service_request_latency_bucket[1m])) by (le))

docker run --name some-redis -d redis
redis-cli
127.0.0.1:6379> ping
PONG
127.0.0.1:6379> select 0
OK
127.0.0.1:6379> keys *
(empty array)

SERVICE_LOG_LEVEL=debug;SERVICE_DB_PASSWD=123456;
# SERVICE_DB_DSN=root:123456@tcp(127.0.0.1:3306)/tokenizer



cloud native apps usually handle multiple types of configurations:
1. environment variables - defaults in code 
2. secrets injected in env vars - not committed to version control
3. configuration files