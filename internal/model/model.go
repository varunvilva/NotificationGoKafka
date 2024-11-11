package model

import "time"

// Recipient struct
type Recipient struct {
	PhoneNumber  string `json:"phone_number,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	DeviceToken  string `json:"device_token,omitempty"`
	Name         string `json:"name,omitempty"`
	Platform     string `json:"platform,omitempty"`
}

// Message struct
type Message struct {
	Subject     string                 `json:"subject,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Content     string                 `json:"content,omitempty"`
	HTML        string                 `json:"html,omitempty"`
	Text        string                 `json:"text,omitempty"`
	Attachments []string               `json:"attachments,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// Sender struct
type Sender struct {
	EmailAddress string `json:"email_address,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Name         string `json:"name,omitempty"`
	SenderID     string `json:"sender_id,omitempty"`
}

// NotificationPayload struct
type NotificationPayload struct {
	Type         string     `json:"type"`
	ServiceName  string     `json:"service_name"`
	Recipient    Recipient  `json:"recipient"`
	Message      Message    `json:"message"`
	Sender       Sender     `json:"sender"`
	RequestTime  time.Time  `json:"request_time,omitempty"`
}

type Test struct{
	Message string `json:"message"`
	Type string `json:"type"`
}
