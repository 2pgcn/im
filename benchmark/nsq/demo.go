package main

import (
	"errors"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type myMessageHandler struct{}

// HandleMessage implements the Handler interface.
func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}
	var processMessage = func([]byte) error {
		fmt.Println("received message: ", string(m.Body))
		if rand.Intn(10)%2 == 0 {
			m.Finish()
			return nil
		}
		m.Requeue(100 * time.Millisecond)
		//m.Finish()
		return errors.New("test111")
	}
	// do whatever actual message processing is desired
	err := processMessage(m.Body)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}

func main() {
	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("topic", "channel", config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	consumer.AddHandler(&myMessageHandler{})
	consumer.ChangeMaxInFlight(20)

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQD("192.168.31.49:4150")
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer.Stop()
}
