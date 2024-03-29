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
# Cache get
GET         /api/v1/cache/{key}
# Cache post
POST        /api/v1/cache/{key} # with plain text body for value
# HealthCheck
GET         /health
# Metrics
GET         /metrics
```

Please import [go-mem-cache.postman_collection.json](https://github.com/abhi195/go-mem-cache/blob/master/go-mem-cache.postman_collection.json) for sample requests.

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

## Future work

- Extend the cache to support more operations.
- Swap custom cache implementation with widely used cache implementations like [github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache) or [github.com/allegro/bigcache](https://github.com/allegro/bigcache) which provides richer functionalities.
