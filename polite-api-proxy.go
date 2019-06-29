package main

func main() {
	var proxyServer ProxyServer
	proxyServer.Setup()
	proxyServer.Run()
	proxyServer.Close()
}
