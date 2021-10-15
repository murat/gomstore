# gomstore

![Status](https://github.com/murat/gomstore/actions/workflows/test.yml/badge.svg)

## Description

**gomstore** is a rest API written in golang that has endpoints to manage an in-memory store.

## Installation & Usage

gomstore developed by using only standard library. Build and run is very simple.

### Locally

```shell
git clone git@github.com:murat/gomstore.git && cd gomstore

go build -o ./bin/gomstore # or make build

go run .
# or
./bin/gomstore
```

### Docker

```shell
git clone git@github.com:murat/gomstore.git && cd gomstore

docker build . -t gomstore
docker run -v /tmp/gomstore.json:/tmp/gomstore.json -p 8080:8080 gomstore
```

gomstore is going to start a http server that listens **8080** port. Click the following Postman button to see api doc.

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/1516159-dae417b8-1ff3-4be7-91b3-6566a6897dfa?action=collection%2Ffork&collection-url=entityId%3D1516159-dae417b8-1ff3-4be7-91b3-6566a6897dfa%26entityType%3Dcollection%26workspaceId%3Da1934f8d-ec93-427f-883e-7aaf6d8f6790)

## Specs

gomstore api talks with memory directly, but it backs up periodically to **/tmp/gomstore.json** file. And it loads the
data to memory when running if the backup file already exists.

## References

- [Different approaches to HTTP routing in Go](https://benhoyt.com/writings/go-routing/)

## Contribute

All PR’s and issues are welcome!