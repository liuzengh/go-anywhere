package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"go-anywhere.com/rpc/server"
)

func main() {
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
