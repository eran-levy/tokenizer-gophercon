# tokenizer-gophercon

Setup pre-requisite - 
docker
microk8s

docker start jaeger
UI: http://localhost:16686/search
apps/prometheus-2.10.0.darwin-amd64 ./prometheus --config.file=prometheus.yml --log.level=error
UI: http://localhost:9090/

cloud native apps usually handle multiple types of configurations:
1. environment variables - defaults in code 
2. secrets injected in env vars - not committed to version control
3. configuration files