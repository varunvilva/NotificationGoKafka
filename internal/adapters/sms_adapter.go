package adapters

import (
    "fmt"
    "github.com/IBM/sarama"
)

type SMSAdapter struct{}

func (s *SMSAdapter) SendNotification(msg *sarama.ConsumerMessage) error {
    fmt.Println("Sending SMS")
    // Implement actual SMS logic here
    return nil
}
