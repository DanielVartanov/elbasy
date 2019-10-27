package main

import (
	"log"
)

func main() {
	elbasy := newProxy(8443)

	defer func() {
		err := elbasy.close()
		if err != nil {
			log.Fatalf("elbasy.close(): %v", err)
		}
	}()

	err := elbasy.run()
	if err != nil {
		log.Fatalf("elbasy.run(): %v", err)
	}
}
