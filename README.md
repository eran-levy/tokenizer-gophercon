## Tokenizer Gophercon sample service
### Gophercon 2021 Workshop
Setup pre-requisite -
* Go 1.15+
* Docker desktop + Kubernetes enabled:
    * Mac- https://docs.docker.com/docker-for-mac/#kubernetes
    * Windows- https://docs.docker.com/docker-for-windows/#kubernetes
    * Linux- https://docs.docker.com/engine/install/ and Minikube installation guide: https://minikube.sigs.k8s.io/docs/start/ (convenience scripts: https://get.docker.com/)
* Skaffold https://skaffold.dev/docs/install/
* Helm https://helm.sh/docs/intro/install/
* kubectl https://kubernetes.io/docs/tasks/tools/install-kubectl/

Make sure you can setup and start (scripts/):
1. MySQL docker container (setup-mysql-server.sh)
2. Redis docker container (setup-redis.sh)
3. Jaeger all-in-one docker container (setup-jaeger-all-in-one.sh)
4. Prometheus local (setup-prometheus.sh)
5. Protoc compiler and go plugins (generate-proto.sh)
6. Execute the sql script to create the database and a table (create-tables.sql)

Useful Web UIs:
* Jaeger UI: http://localhost:16686/search
* Prometheus UI: http://localhost:9090/
Redis CLI:

`redis-cli
127.0.0.1:6379> ping
PONG
127.0.0.1:6379> select 0
OK
127.0.0.1:6379> keys *
(empty array)`

Service environment variables override (DSN env var: SERVICE_DB_DSN=root:123456@tcp(127.0.0.1:3306)/tokenizer):

`SERVICE_LOG_LEVEL=debug;SERVICE_DB_PASSWD=123456;`


