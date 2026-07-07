package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)
type NotificationEvent struct{
	Message string
}
type Notifier struct{}

func main(){
	notifier := new(Notifier)
	rpc.Register(notifier)

	listener, err := net.Listen("tcp",":8082")
	if err !=nil{
		log.Fatal("Błąd nasłuchu RPC: ", err)
	}

	fmt.Println("Mikroserwis Powiadomień - :8082")
	for{
		conn, err := listener.Accept()
		if err != nil{
			continue
		}
		go rpc.ServeConn(conn)
	}
}