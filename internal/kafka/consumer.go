package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"github.com/IBM/sarama"
    "NotificationMicroservice/internal/adapters"
    "log"
)

// List of topics to consume
var topics = []string{"email_notifications", "sms_notifications", "push_notifications"}
var emailAdapter adapters.EmailAdapter
var pushAdapter adapters.PushAdapter
var smsAdapter adapters.SMSAdapter

type NotificationService struct {
    emailAdapter adapters.Adapter
    pushAdapter  adapters.Adapter
    smsAdapter   adapters.Adapter
}

func NewNotificationService(email adapters.Adapter, push adapters.Adapter, sms adapters.Adapter) *NotificationService {
    return &NotificationService{
        emailAdapter: email,
        pushAdapter:  push,
        smsAdapter:   sms,
    }
}

func main() {
	brokers := []string{"localhost:29092"}
	
	// Connect to Kafka consumer
	consumer, err := connectConsumer(brokers)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println("Error closing consumer:", err)
		}
	}()

	// Set up channels to handle OS signals and to synchronize message processing
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	doneCh := make(chan struct{})

	// Wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Loop through each topic and consume all its partitions concurrently
	for _, topic := range topics {
		partitions, err := consumer.Partitions(topic)
		if err != nil {
			panic(err)
		}

		for _, partition := range partitions {
			// Increase wait group counter for each partition consumer
			wg.Add(1)
			
			// Start a goroutine to consume messages from each partition
			go func(topic string, partition int32) {
				defer wg.Done() // Signal completion of this goroutine
				consumePartition(consumer, topic, partition, doneCh)
			}(topic, partition)
		}
	}

	// Wait for interruption signal to shut down
	<-sigchan
	close(doneCh) // Signal all consumers to stop
	wg.Wait()     // Wait for all goroutines to complete
	fmt.Println("Consumer shutdown complete.")
}
func SendNotification(message *sarama.ConsumerMessage) {
    fmt.Printf("Received message | Topic: %s | Partition: %d | Offset: %d | Message: %s\n",message .Topic, message .Partition, message.Offset, string(message.Value))
    var err error
    switch message.Topic {
    case "email_notifications":
        err = emailAdapter.SendNotification(message)
    case "push_notification":
        err = pushAdapter.SendNotification(message)
    case "sms_notifications":
        err = smsAdapter.SendNotification(message)
    default:
        log.Printf("unsupported notification type: %s", message.Topic)
        return
    }

    if err != nil {
        log.Printf("error sending notification: %v", err)
        // Logic to handle failed notifications, e.g., send to failed_notifications topic
    }
}

// consumePartition consumes messages from a specific topic partition
func consumePartition(consumer sarama.Consumer, topic string, partition int32, doneCh <-chan struct{}) {
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("Failed to start consumer for topic %s partition %d: %v\n", topic, partition, err)
		return
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			fmt.Println("Error closing partition consumer:", err)
		}
	}()

	fmt.Printf("Started consuming topic %s partition %d\n", topic, partition)

	for {
		select {
		case err := <-partitionConsumer.Errors():
			fmt.Println("Error:", err)
		case msg := <-partitionConsumer.Messages():
			// Process each message asynchronously
			go SendNotification(msg)
		case <-doneCh:
			// Exit when doneCh is closed
			return
		}
	}
}



// handleMessage processes messages based on the topic they were received from
// func handleMessage(msg *sarama.ConsumerMessage) {
// 	fmt.Printf("Received message | Topic: %s | Partition: %d | Offset: %d | Message: %s\n",
// 		msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

// 	// Route the message to the appropriate function based on the topic
// 	switch msg.Topic {
// 	case "email_notifications":
// 		handleEmailNotification(msg)
// 	case "sms_notifications":
// 		handleSMSNotification(msg)
// 	case "push_notifications":
// 		handlePushNotification(msg)
// 	default:
// 		fmt.Printf("Unknown topic: %s\n", msg.Topic)
// 	}
// }

// Specific handlers for each topic
// func handleTopic1Message(msg *sarama.ConsumerMessage) {
// 	fmt.Printf("Handling message for topic1: %s\n", string(msg.Value))
// 	// Add business logic for topic1 messages here
// }

// func handleTopic2Message(msg *sarama.ConsumerMessage) {
// 	fmt.Printf("Handling message for topic2: %s\n", string(msg.Value))
// 	// Add business logic for topic2 messages here
// }

// func handleTopic3Message(msg *sarama.ConsumerMessage) {
// 	fmt.Printf("Handling message for topic3: %s\n", string(msg.Value))
// 	// Add business logic for topic3 messages here
// }

// Function to connect to the Kafka consumer
func connectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
