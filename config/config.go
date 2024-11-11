package config

import (
    "gopkg.in/yaml.v2"
    "os"
)

type AdapterConfig struct {
    Name        string            `yaml:"name"`
    Priority    int               `yaml:"priority"`
    MaxRetries  int               `yaml:"max_retries"`
    Credentials map[string]string `yaml:"credentials"`
}

type NotificationPolicy struct {
    SMS struct {
        Adapters []AdapterConfig `yaml:"adapters"`
    } `yaml:"sms"`
    Email struct {
        Adapters []AdapterConfig `yaml:"adapters"`
    } `yaml:"email"`
    Push struct {
        Adapters []AdapterConfig `yaml:"adapters"`
    } `yaml:"push"`
}

type KafkaConfig struct {
    Brokers []string `yaml:"brokers"`
    Topics  struct {
        Email string `yaml:"email"`
        SMS   string `yaml:"sms"`
        Push  string `yaml:"push"`
		Failed string `yaml:"failed"`
    } `yaml:"topics"`
}

type Config struct {
    NotificationPolicies NotificationPolicy `yaml:"notification_policies"`
    Kafka                KafkaConfig        `yaml:"kafka"`
}

// LoadConfig loads and parses the YAML configuration file
func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}
