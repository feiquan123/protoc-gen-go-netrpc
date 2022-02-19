package main

import (
	"log"
	"net/rpc"

	"github.com/feiquan123/protoc-gen-go-netrpc/example/api"
)

type HelloService struct{}

func (p *HelloService) Hello(request *api.String, replay *api.String) error {
	replay.Value = "hello:" + request.GetValue()
	return nil
}

func main() {
	run, err := api.GenRunHelloService(
		rpc.DefaultServer,
		new(HelloService),
		"tcp", ":1234",
		api.EncodeJson,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(run())
}
