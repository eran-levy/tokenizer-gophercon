# tokenizer-gophercon

Setup pre-requisite - 
docker
microk8s

cloud native apps usually handle multiple types of configurations:
1. environment variables - defaults in code 
2. secrets injected in env vars - not committed to version control
3. configuration files