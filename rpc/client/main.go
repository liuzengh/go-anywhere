package main

import (
	"fmt"
	"log"
	"net/rpc"

	"go-anywhere.com/rpc/server"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := &server.Args{A: 7, B: 8}
	var reply int
	// err = client.Call("Arith.Multiply", args, &reply)
	call := client.Go("Arith.Multiply", args, &reply, nil)
	replyCall := <-call.Done
	if replyCall.Error != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)
}
