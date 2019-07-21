package main

import (
	"log"
	"./proxy_server"
)

func main() {
	var proxyServer proxy_server.ProxyServer
	err := proxyServer.Run()
	if err != nil {
		log.Fatalf("Error ProxyServer.Run(): %v", err)
	}

	err = proxyServer.Close()
	if err != nil {
		log.Fatalf("Error ProxyServer.Close(): %v", err)
	}
}
