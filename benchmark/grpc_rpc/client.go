package main

import (
	"context"
	"flag"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-plugins/transport/grpc"

	"github.com/hb-go/micro/benchmark/service"
)

var concurrency = flag.Int("c", 1, "concurrency")
var total = flag.Int("n", 1, "total requests for all clients")
var contentType = flag.String("ct", "application/protobuf", "content type")

func main() {
	flag.Parse()
	serviceName := "hello.grpc.rpc"
	service.ClientWith(
		serviceName,
		func() client.Client {
			cache := selector.NewSelector(func(o *selector.Options) {
				o.Context = context.WithValue(o.Context, "selector_ttl", time.Minute*3)
			})
			return client.NewClient(
				client.Transport(grpc.NewTransport()),
				client.ContentType(*contentType),
				client.Selector(cache),
				client.Retries(1),
				client.PoolSize(*concurrency*2),
				client.RequestTimeout(time.Millisecond*100),
				// client.Wrap(breaker.NewClientWrapper()),
				// client.Wrap(ratelimit.NewClientWrapper(10000)),
			)
		},
		*concurrency,
		*total,
	)

}
