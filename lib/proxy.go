package lib

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSProxy struct {
	svc    *sqs.SQS
	url    string
	target string
}

func NewSQSProxy(url, target string) (*SQSProxy, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &SQSProxy{
		svc:    sqs.New(sess),
		url:    url,
		target: target,
	}, nil
}

func (s *SQSProxy) Run() error {
	errorCount := 0
	for {
		err := s.getMessage()
		if err != nil {
			errorCount++
		}
		if errorCount > 10 {
			return err
		}
	}
}

// from https://qiita.com/sudix/items/215d7ffb65e89187c1f1
func (s *SQSProxy) getMessage() error {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.url),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
	}
	resp, err := s.svc.ReceiveMessage(params)
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
			err := s.HttpPost(*msg.Body)
			if err != nil {
				// ignore error, should I use retry queue?
				log.Fatal(err)
			}
			if err := s.DeleteMessage(msg); err != nil {
				// rare case
				log.Fatal(err)
			}
		}(m)
	}

	wg.Wait()
	return nil
}

func (s *SQSProxy) DeleteMessage(msg *sqs.Message) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.url),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	if _, err := s.svc.DeleteMessage(params); err != nil {
		return err
	}
	return nil
}

func (s *SQSProxy) HttpPost(value string) error {
	if len(value) < 8 {
		log.Println("Invalid message")
		return nil
	}

	values := url.Values{}
	buf, err := url.QueryUnescape(value[8:])
	if err != nil {
		return err
	}
	values.Add("payload", buf)

	req, err := http.NewRequest(
		"POST",
		s.target,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
