# go-mem-cache

A simple in-memory cache with HTTP REST interface.

## Cache

### usage

```go
package main

import (
	"fmt"
	"time"

	memcache "github.com/abhi195/go-mem-cache"
)

func main() {

	// creating cache with default 30 minutes of TTL which
	// evicts expired entries every 60 minutes
	mc := memcache.New()

	// alternatively you can use NewWith() to
	// create cache with required values
	nmc := memcache.NewWith(1*time.Minute, 2*time.Minute)

	// set key, value
	mc.Set("foo", "bar")

	// get key
	v, exist := mc.Get("foo")
	if exist {
		fmt.Println(v)
	}
}
```

## REST Interface

### APIs

```bash
GET         /api/v1/cache/{key}
POST        /api/v1/cache/{key}
```

Please refer to [go-mem-cache.postman_collection.json](https://github.com/abhi195/go-mem-cache/blob/master/go-mem-cache.postman_collection.json) file for sample requests.

### Starting HTTP Server

- Natively

```
go build -o ./out/cache-server server/*
./out/cache-server
```

- Using provided docker image

```
docker pull abhi195/go-mem-cache
docker run --rm -i -p 8080:8080 abhi195/go-mem-cache
```

- Building docker image and running

```
docker build -t go-mem-cache .
docker run --rm -i -p 8080:8080 go-mem-cache
```