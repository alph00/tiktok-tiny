package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
	"github.com/cloudwego/kitex/client"
)

func main() {
	client, err := auth.NewClient("hello", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		req := &auth.Request{Message: "my request"}
		resp, err := client.Register(context.Background(), req)

		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)
		addreq := &api.AddRequest{First: 512, Second: 512}
		addresp, err := client.Add(context.Background(), addreq)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(addresp)
		time.Sleep(time.Second)

	}
}