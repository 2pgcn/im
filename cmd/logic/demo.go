package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func demo() {
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0
	config.Producer.Return.Successes = false
	config.Producer.RequiredAcks = sarama.NoResponse
	producer, err := sarama.NewAsyncProducer([]string{"192.168.1.123:9092"}, config)
	if err != nil {
		panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	var (
		wg                        sync.WaitGroup
		successes, producerErrors int
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Successes() {
			log.Println(err)
			producerErrors++
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Successes() {
			log.Println(err)
			producerErrors++
		}
	}()
	for i := 0; i <= 10000; i++ {
		go func() {
			for i := 0; i < 10000; i++ {
				producer.Input() <- &sarama.ProducerMessage{
					Topic:     "test",
					Value:     sarama.StringEncoder("hello world!"),
					Partition: int32(i % 9),
				}
			}
		}()
	}
For:
	for {
		select {
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			break For
		default:
			time.Sleep(time.Second * 2)
			fmt.Printf("success: %d, error: %d\n", successes, producerErrors)
		}
	}
	wg.Wait()

}
