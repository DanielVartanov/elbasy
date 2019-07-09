package main

func main() {
	var proxyServer ProxyServer
	proxyServer.Run()
	proxyServer.Close()
}
