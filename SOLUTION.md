# GitHub Gist Proxy API Service

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Go](https://go.dev/dl/) version 1.27

## Run tests

```bash
make test
```

## Build docker image

```bash
docker build -t gist-service:latest .
```

## Run locally

```bash
make run
```

## Test for 'octocat' user

```bash
curl http://localhost:8080/octocat
```

## Links

[List gists for a user - api reference](https://docs.github.com/en/rest/gists/gists?apiVersion=2022-11-28#list-gists-for-a-user)

[go-github client library](https://github.com/google/go-github/tree/master)

[go-github client library - list gists for a user](https://github.com/google/go-github/blob/19ae9a1abca8f30c71a8e113618f60bcadddcddb/github/gists.go#L138)
