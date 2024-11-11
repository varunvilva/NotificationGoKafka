package adapters

import (
    "fmt"
    "github.com/IBM/sarama"
)

type PushAdapter struct{}

func (p *PushAdapter) SendNotification(msg *sarama.ConsumerMessage) error {
    fmt.Println("Sending push notification")
    // Implement actual push logic here
    return nil
}
