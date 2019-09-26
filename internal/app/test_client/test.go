package main

import (
	"context"
	"crawler/internal/svc/lib"
	"crawler/internal/svc/rpc"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("127.0.0.1:7981", opts...)
	if err != nil {
		log.Fatal("[error] connect error " + err.Error())
	}
	defer conn.Close()

	client := rpc.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()

	var result *rpc.Result
	result, err = client.Craw(ctx, &rpc.Job{Name: lib.SITE_WEIBO})
	fmt.Println(result)
}
