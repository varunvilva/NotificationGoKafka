package adapters

import ("fmt"
        "github.com/IBM/sarama"
        )
type EmailAdapter struct{}

func (e *EmailAdapter) SendNotification(msg *sarama.ConsumerMessage) error {
    fmt.Println("Sending email")
    // Implement actual email logic here
    return nil
}
// func handleTopic1Message(msg *sarama.ConsumerMessage) {
//     // 	fmt.Printf("Handling message for topic1: %s\n", string(msg.Value))
//     // 	// Add business logic for topic1 messages here
//     }