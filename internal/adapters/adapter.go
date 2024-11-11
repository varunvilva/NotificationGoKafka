package adapters

import "github.com/IBM/sarama"

type Adapter interface {
    SendNotification(msg *sarama.ConsumerMessage) error
}