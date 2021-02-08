#!/usr/bin/make -f

docker-build:
	docker build -t tokenizer-gophercon:test .

docker-run:
	docker run -p 8080:8080 -p 7070:7070 --name mytest1 tokenizer-gophercon:test

docker-cleanup:
	docker rm mytest1

skaffold-dev:
	skaffold dev

helm-template:
	 helm template helm/tokenizer-gophercon --validate
