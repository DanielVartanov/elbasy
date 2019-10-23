package main

import (
	"log"
)

func main() {
	proxy := newProxy()

	defer func() {
		err := proxy.close()
		if err != nil {
			log.Fatalf("proxy.close(): %v", err)
		}
	}()

	err := proxy.run()
	if err != nil {
		log.Fatalf("proxy.run(): %v", err)
	}
}
