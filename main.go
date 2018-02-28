package main

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	svc := sqs.New(&aws.Config{Region: aws.String(AWS_REGION)})
	url := ""
	GetMessage(svc, url)
}

func GetMessage(svc *sqs.SQS, url string) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(url),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
	}
	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		return err
	}

	fmt.Printf("messages count: %d\n", len(resp.Messages))

	if len(resp.Messages) == 0 {
		fmt.Println("empty queue.")
		return nil
	}

	var wg sync.WaitGroup
	for _, m := range resp.Messages {
		wg.Add(1)
		go func(msg *sqs.Message) {
			defer wg.Done()
			fmt.Println(msg.Body)
			if err := DeleteMessage(msg); err != nil {
				fmt.Println(err)
			}
		}(m)
	}

	wg.Wait()

	return nil
}
