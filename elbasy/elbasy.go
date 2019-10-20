package main

import (
	"log"
)

func main() {
	proxy := newProxy()

	defer func() {
		err := proxy.close()
		if err != nil {
			log.Fatalf("Error proxy.Close(): %v", err)
		}
	}()

	err := proxy.run()
	if err != nil {
		log.Fatalf("Error proxy.Run(): %v", err)
	}
}
