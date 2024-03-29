package tasks

import "encoding/json"

const (
	EmailInvite       = "task_email_invite"
	EmailVerification = "task_email_verification"
)

type PublishMessage struct {
	Task    string `json:"task"`
	Payload any    `json:"payload"`
}

type ConsumeMessage struct {
	Task    string          `json:"task"`
	Payload json.RawMessage `json:"payload"`
}
