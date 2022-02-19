package main

import (
	"fmt"
	"log"

	"github.com/feiquan123/protoc-gen-go-netrpc/example/api"
)

func main() {
	client, err := api.DialHelloService("tcp", "localhost:1234", api.EncodeJson)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	in := api.String{
		Value: "world",
	}
	out := api.String{}
	err = client.Hello(&in, &out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.GetValue())
}
