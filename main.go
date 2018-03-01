package main

import (
	"os"

	"log"

	"github.com/juntaki/firestarter-sqs-proxy/lib"
)

func main() {
	url := os.Getenv("SQS_URL")
	target := os.Getenv("POST_URL")
	proxy, err := lib.NewSQSProxy(url, target)
	if err != nil {
		log.Fatal("Failed to make SQS proxy", err)
	}
	err = proxy.Run() // should not return
	log.Fatal("Failed", err)
}
