package main

import (
	"log"
)

func main() {
	var proxyServer ProxyServer
	err := proxyServer.Run()
	if err != nil {
		log.Fatalf("Error ProxyServer.Run(): %v", err)
	}

	err = proxyServer.Close()
	if err != nil {
		log.Fatalf("Error ProxyServer.Close(): %v", err)
	}
}
